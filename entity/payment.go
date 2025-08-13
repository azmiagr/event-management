package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	PaymentID      uuid.UUID `json:"payment_id" gorm:"type:varchar(36);primaryKey"`
	RegistrationID uuid.UUID `json:"registration_id"`
	Amount         float64   `json:"amount" gorm:"type:float(8,2)"`
	Method         string    `json:"method" gorm:"type:varchar(20)"`
	Status         string    `json:"status" gorm:"type:enum('pending','paid','failed')"`
	PaidAt         time.Time `json:"paid_at" gorm:"type:datetime;default:null"`
}
