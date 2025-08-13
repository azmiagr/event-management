package entity

import (
	"time"

	"github.com/google/uuid"
)

type Registration struct {
	RegistrationID uuid.UUID `json:"registration_id" gorm:"type:varchar(36);primaryKey"`
	EventID        uuid.UUID `json:"event_id"`
	UserID         uuid.UUID `json:"user_id"`
	TicketCode     string    `json:"ticket_code" gorm:"type:varchar(6)"`
	Status         string    `json:"status" gorm:"type:enum('pending', 'approved', 'rejected');default:'pending'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Payment Payment `json:"payment" gorm:"foreignKey:RegistrationID"`
}
