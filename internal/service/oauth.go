package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/bcrypt"
	"event-management/pkg/config"
	"event-management/pkg/database/mariadb"
	"event-management/pkg/jwt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOAuthService interface {
	GetGoogleLoginURL() (string, string, error)
	HandleGoogleCallback(code, state, savedState string) (*model.OAuthLoginResponse, error)
	GetUserInfoFromGoogle(accessToken string) (*model.GoogleUserInfo, error)
}

type OAuthService struct {
	db             *gorm.DB
	UserRepository repository.IUserRepository
	oauth          *config.OAuthConfig
	bcrypt         bcrypt.Interface
	jwt            jwt.Interface
}

func NewOAuthService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface, jwt jwt.Interface, oauth *config.OAuthConfig) IOAuthService {
	return &OAuthService{
		db:             mariadb.Connection,
		UserRepository: userRepository,
		oauth:          oauth,
		bcrypt:         bcrypt,
		jwt:            jwt,
	}
}

func (s *OAuthService) GetGoogleLoginURL() (string, string, error) {
	state, err := s.generateState()
	if err != nil {
		return "", "", err
	}

	url := s.oauth.GoogleConfig.AuthCodeURL(state)
	return url, state, nil
}

func (s *OAuthService) HandleGoogleCallback(code, state, savedState string) (*model.OAuthLoginResponse, error) {
	if state != savedState {
		return nil, errors.New("invalid state parameter")
	}

	token, err := s.oauth.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code for token")
	}

	googleUser, err := s.GetUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.findOrCreateUser(googleUser)
	if err != nil {
		return nil, err
	}

	jwtToken, err := s.jwt.CreateJWTToken(user.UserID, false)
	if err != nil {
		return nil, errors.New("failed to create JWT token")
	}

	return &model.OAuthLoginResponse{
		Token: jwtToken,
	}, nil
}

func (s *OAuthService) GetUserInfoFromGoogle(accessToken string) (*model.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo model.GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *OAuthService) findOrCreateUser(googleUser *model.GoogleUserInfo) (*entity.User, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	user, err := s.UserRepository.GetUser(tx, model.GetUserParam{
		GoogleID: &googleUser.ID,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user != nil {
		user.Name = googleUser.Name
		user.Picture = &googleUser.Picture
		_, err = s.UserRepository.UpdateUser(tx, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	existingUser, err := s.UserRepository.GetUser(tx, model.GetUserParam{
		Email: googleUser.Email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		existingUser.GoogleID = &googleUser.ID
		existingUser.Picture = &googleUser.Picture

		_, err = s.UserRepository.UpdateUser(tx, existingUser)
		if err != nil {
			return nil, err
		}

		return existingUser, nil
	}

	newUser := &entity.User{
		UserID:        uuid.New(),
		GoogleID:      &googleUser.ID,
		Name:          googleUser.Name,
		Picture:       &googleUser.Picture,
		Email:         googleUser.Email,
		StatusAccount: "active",
		RoleID:        2,
	}

	createdUser, err := s.UserRepository.CreateUser(tx, newUser)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return createdUser, nil

}

func (s *OAuthService) generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	return state, nil
}
