package service

import (
	"event-management/internal/repository"
	"event-management/pkg/bcrypt"
	"event-management/pkg/config"
	"event-management/pkg/jwt"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Service struct {
	UserService         IUserService
	OAuthService        IOAuthService
	OtpService          IOtpService
	RegistrationService IRegistrationService
	EventService        IEventService
	PaymentService      IPaymentService
}

func NewService(repo *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, oauth *config.OAuthConfig, snapClient snap.Client, coreAPIClient coreapi.Client) *Service {
	paymentService := NewPaymentService(repo.PaymentRepository, snapClient, coreAPIClient)

	return &Service{
		UserService:         NewUserService(repo.UserRepository, repo.OtpRepository, bcrypt, jwtAuth),
		OAuthService:        NewOAuthService(repo.UserRepository, bcrypt, jwtAuth, oauth),
		OtpService:          NewOtpService(repo.OtpRepository, repo.UserRepository),
		RegistrationService: NewRegistrationService(repo.RegistrationRepository, repo.EventRepository, paymentService),
		EventService:        NewEventService(repo.EventRepository),
		PaymentService:      paymentService,
	}
}
