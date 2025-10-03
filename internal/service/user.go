package service

import (
	"errors"
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/bcrypt"
	"event-management/pkg/database/mariadb"
	"event-management/pkg/jwt"

	"gorm.io/gorm"
)

type IUserService interface {
	Register(param model.UserRegisterParam) error
	Login(param model.UserLoginParam) (*model.UserLoginResponse, error)
}

type UserService struct {
	db             *gorm.DB
	bcrypt         bcrypt.Interface
	jwt            jwt.Interface
	UserRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface, jwt jwt.Interface) IUserService {
	return &UserService{
		db:             mariadb.Connection,
		bcrypt:         bcrypt,
		jwt:            jwt,
		UserRepository: userRepository,
	}
}

func (s *UserService) Register(param model.UserRegisterParam) error {
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

	user := &entity.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: &hash,
	}

	_, err = s.UserRepository.CreateUser(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(param model.UserLoginParam) (*model.UserLoginResponse, error) {
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
