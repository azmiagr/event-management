package repository

import (
	"event-management/entity"

	"gorm.io/gorm"
)

type IPaymentRepository interface {
	CreatePayment(tx *gorm.DB, payment *entity.Payment) (*entity.Payment, error)
	GetPaymentByID(tx *gorm.DB, id string) (*entity.Payment, error)
	UpdatePayment(tx *gorm.DB, payment *entity.Payment) (*entity.Payment, error)
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(tx *gorm.DB, payment *entity.Payment) (*entity.Payment, error) {
	err := tx.Debug().Create(&payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) GetPaymentByID(tx *gorm.DB, id string) (*entity.Payment, error) {
	var payment entity.Payment
	err := tx.Debug().Where("order_id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) UpdatePayment(tx *gorm.DB, payment *entity.Payment) (*entity.Payment, error) {
	err := tx.Debug().Save(&payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil
}
