package service

import (
	"event-management/internal/repository"
	"event-management/pkg/bcrypt"
	"event-management/pkg/config"
	"event-management/pkg/jwt"
)

type Service struct {
	UserService  IUserService
	OAuthService IOAuthService
	OtpService   IOtpService
}

func NewService(repo *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, oauth *config.OAuthConfig) *Service {
	return &Service{
		UserService:  NewUserService(repo.UserRepository, repo.OtpRepository, bcrypt, jwtAuth),
		OAuthService: NewOAuthService(repo.UserRepository, bcrypt, jwtAuth, oauth),
		OtpService:   NewOtpService(repo.OtpRepository, repo.UserRepository),
	}
}
