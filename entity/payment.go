package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	OrderID        string    `json:"order_id" gorm:"type:varchar(255);primaryKey"`
	RegistrationID uuid.UUID `json:"registration_id"`
	Amount         float64   `json:"amount" gorm:"type:float(8,2)"`
	Status         string    `json:"status" gorm:"type:varchar(30)"`
	SnapURL        string    `json:"snap_url" gorm:"type:varchar(255)"`
	PaymentType    string    `json:"payment_type" gorm:"type:varchar(50)"`
	PaidAt         time.Time `json:"paid_at" gorm:"type:datetime;default:null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
