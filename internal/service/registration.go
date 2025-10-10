package service

import (
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/database/mariadb"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRegistrationService interface {
	RegisterEvent(userID, eventID uuid.UUID) (*model.CreatePaymentResponse, error)
	UpdateStatusAfterPayment(tx *gorm.DB, registrationID uuid.UUID, status string) error
}

type RegistrationService struct {
	db                     *gorm.DB
	RegistrationRepository repository.IRegistrationRepository
	EventRepository        repository.IEventRepository
	PaymentService         IPaymentService
}

func NewRegistrationService(registrationRepository repository.IRegistrationRepository, eventRepository repository.IEventRepository, paymentService IPaymentService) IRegistrationService {
	return &RegistrationService{
		db:                     mariadb.Connection,
		RegistrationRepository: registrationRepository,
		EventRepository:        eventRepository,
		PaymentService:         paymentService,
	}
}

func (s *RegistrationService) RegisterEvent(userID, eventID uuid.UUID) (*model.CreatePaymentResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	event, err := s.EventRepository.GetEvent(tx, model.EventParam{
		EventID: eventID,
	})
	if err != nil {
		return nil, err
	}

	registrationID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	registration := &entity.Registration{
		RegistrationID: registrationID,
		EventID:        eventID,
		UserID:         userID,
		TicketCode:     generateTicketCode(),
		Status:         "pending",
	}

	_, err = s.RegistrationRepository.CreateRegistration(tx, registration)
	if err != nil {
		return nil, err
	}

	snapResp, err := s.PaymentService.CreatePayment(tx, registration.RegistrationID, event.Price)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

func (s *RegistrationService) UpdateStatusAfterPayment(tx *gorm.DB, registrationID uuid.UUID, status string) error {
	registration, err := s.RegistrationRepository.GetRegistration(tx, model.RegistrationParam{
		RegistrationID: registrationID,
	})
	if err != nil {
		return err
	}

	registration.Status = status

	_, err = s.RegistrationRepository.UpdateRegistration(tx, registration)
	if err != nil {
		return err
	}

	return nil
}

func generateTicketCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}
