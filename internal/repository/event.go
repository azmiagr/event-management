package repository

import (
	"event-management/entity"
	"event-management/model"

	"gorm.io/gorm"
)

type IEventRepository interface {
	CreateEvent(tx *gorm.DB, event *entity.Event) (*entity.Event, error)
	GetEvent(tx *gorm.DB, param model.EventParam) (*entity.Event, error)
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) IEventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(tx *gorm.DB, event *entity.Event) (*entity.Event, error) {
	err := tx.Create(&event).Error
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) GetEvent(tx *gorm.DB, param model.EventParam) (*entity.Event, error) {
	var event entity.Event

	err := tx.Debug().Where(&param).First(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}
