package service

import (
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/database/mariadb"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IEventService interface {
	CreateEvent(userID uuid.UUID, param *model.CreateEventParam) (*model.CreateEventResponse, error)
}

type EventService struct {
	db              *gorm.DB
	EventRepository repository.IEventRepository
}

func NewEventService(eventRepo repository.IEventRepository) IEventService {
	return &EventService{
		db:              mariadb.Connection,
		EventRepository: eventRepo,
	}
}

func (s *EventService) CreateEvent(userID uuid.UUID, param *model.CreateEventParam) (*model.CreateEventResponse, error) {
	var response *model.CreateEventResponse

	tx := s.db.Begin()
	defer tx.Rollback()

	eventID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	event := &entity.Event{
		EventID:      eventID,
		UserID:       userID,
		Title:        param.Title,
		Description:  param.Description,
		Category:     param.Category,
		StartDate:    param.StartDate,
		EndDate:      param.EndDate,
		Location:     param.Location,
		LocationType: param.LocationType,
		Quota:        param.Quota,
		Price:        param.Price,
		Status:       "draft",
	}

	_, err = s.EventRepository.CreateEvent(tx, event)
	if err != nil {
		return nil, err
	}

	response = &model.CreateEventResponse{
		EventID:      event.EventID,
		Title:        event.Title,
		Description:  event.Description,
		Category:     event.Category,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		Location:     event.Location,
		LocationType: event.LocationType,
		Quota:        event.Quota,
		Price:        event.Price,
		Status:       event.Status,
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return response, nil

}
