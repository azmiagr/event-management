package repository

import (
	"event-management/entity"
	"event-management/model"

	"gorm.io/gorm"
)

type IOtpRepository interface {
	GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.Otp, error)
	CreateOtp(tx *gorm.DB, otp *entity.Otp) error
	UpdateOtp(tx *gorm.DB, otp *entity.Otp) error
	DeleteOtp(tx *gorm.DB, otp *entity.Otp) error
}

type OtpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) IOtpRepository {
	return &OtpRepository{
		db: db,
	}
}

func (o *OtpRepository) GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.Otp, error) {
	var otp *entity.Otp
	err := tx.Debug().Where(&param).First(&otp).Error
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (o *OtpRepository) CreateOtp(tx *gorm.DB, otp *entity.Otp) error {
	err := tx.Debug().Create(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) UpdateOtp(tx *gorm.DB, otp *entity.Otp) error {
	err := tx.Debug().Updates(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) DeleteOtp(tx *gorm.DB, otp *entity.Otp) error {
	err := tx.Debug().Delete(otp).Error
	if err != nil {
		return err
	}

	return nil
}
