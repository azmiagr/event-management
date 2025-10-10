package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository         IUserRepository
	OtpRepository          IOtpRepository
	PaymentRepository      IPaymentRepository
	RegistrationRepository IRegistrationRepository
	EventRepository        IEventRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:         NewUserRepository(db),
		OtpRepository:          NewOtpRepository(db),
		PaymentRepository:      NewPaymentRepository(db),
		RegistrationRepository: NewRegistrationRepository(db),
		EventRepository:        NewEventRepository(db),
	}
}
