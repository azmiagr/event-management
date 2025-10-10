package repository

import (
	"event-management/entity"
	"event-management/model"

	"gorm.io/gorm"
)

type IRegistrationRepository interface {
	CreateRegistration(tx *gorm.DB, registration *entity.Registration) (*entity.Registration, error)
	GetRegistration(tx *gorm.DB, param model.RegistrationParam) (*entity.Registration, error)
	UpdateRegistration(tx *gorm.DB, registration *entity.Registration) (*entity.Registration, error)
}

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) IRegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) CreateRegistration(tx *gorm.DB, registration *entity.Registration) (*entity.Registration, error) {
	err := tx.Debug().Create(&registration).Error
	if err != nil {
		return nil, err
	}

	return registration, nil
}

func (r *RegistrationRepository) GetRegistration(tx *gorm.DB, param model.RegistrationParam) (*entity.Registration, error) {
	var registration entity.Registration

	err := tx.Debug().Where(&param).First(&registration).Error
	if err != nil {
		return nil, err
	}

	return &registration, nil
}

func (r *RegistrationRepository) UpdateRegistration(tx *gorm.DB, registration *entity.Registration) (*entity.Registration, error) {
	err := tx.Debug().Save(&registration).Error
	if err != nil {
		return nil, err
	}

	return registration, nil
}
