package service

import (
	"errors"
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/bcrypt"
	"event-management/pkg/database/mariadb"
	"event-management/pkg/jwt"
	"event-management/pkg/mail"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param *model.UserRegisterParam) error
	Login(param *model.UserLoginParam) (*model.UserLoginResponse, error)
	VerifyUser(param model.VerifyUser) error
	GetUser(param model.GetUserParam) (*entity.User, error)
}

type UserService struct {
	db             *gorm.DB
	bcrypt         bcrypt.Interface
	jwt            jwt.Interface
	UserRepository repository.IUserRepository
	OtpRepository  repository.IOtpRepository
}

func NewUserService(userRepository repository.IUserRepository, OtpRepository repository.IOtpRepository, bcrypt bcrypt.Interface, jwt jwt.Interface) IUserService {
	return &UserService{
		db:             mariadb.Connection,
		bcrypt:         bcrypt,
		jwt:            jwt,
		UserRepository: userRepository,
		OtpRepository:  OtpRepository,
	}
}

func (s *UserService) Register(param *model.UserRegisterParam) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	existingUser, err := s.UserRepository.GetUser(tx, model.GetUserParam{
		Email: param.Email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingUser != nil {
		return errors.New("user already exists")
	}

	if param.Password != param.ConfirmPassword {
		return errors.New("password does not match")
	}

	hash, err := s.bcrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user := &entity.User{
		UserID:   id,
		RoleID:   2,
		Name:     param.Name,
		Email:    param.Email,
		Password: &hash,
	}

	_, err = s.UserRepository.CreateUser(tx, user)
	if err != nil {
		return err
	}

	code := mail.GenerateCode()
	otp := &entity.Otp{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   code,
	}

	err = s.OtpRepository.CreateOtp(tx, otp)
	if err != nil {
		return err
	}

	err = mail.SendEmail(user.Email, "OTP Verification", code)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(param *model.UserLoginParam) (*model.UserLoginResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	user, err := s.UserRepository.GetUser(tx, model.GetUserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = s.bcrypt.CompareAndHashPassword(*user.Password, param.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.jwt.CreateJWTToken(user.UserID, false)
	if err != nil {
		return nil, err
	}

	result := &model.UserLoginResponse{
		Token: token,
	}

	return result, nil
}

func (s *UserService) VerifyUser(param model.VerifyUser) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	otp, err := s.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if otp.Code != param.OtpCode {
		return errors.New("invalid otp code")
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return err
	}

	expiredThreshold := time.Now().UTC().Add(-time.Duration(expiredTime) * time.Minute)
	if otp.UpdatedAt.Before(expiredThreshold) {
		return errors.New("otp expired")
	}

	user, err := s.UserRepository.GetUser(tx, model.GetUserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	user.StatusAccount = "active"
	_, err = s.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return err
	}

	err = s.OtpRepository.DeleteOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUser(param model.GetUserParam) (*entity.User, error) {
	return u.UserRepository.GetUser(u.db, param)
}
