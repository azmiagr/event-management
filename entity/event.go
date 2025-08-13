package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID      uuid.UUID `json:"event_id" gorm:"type:varchar(36);primaryKey"`
	UserID       uuid.UUID `json:"user_id"`
	Title        string    `json:"title" gorm:"type:varchar(70)"`
	Description  string    `json:"description" gorm:"type:text"`
	Category     string    `json:"category" gorm:"type:varchar(50)"`
	StartDate    time.Time `json:"start_date" gorm:"type:date"`
	EndDate      time.Time `json:"end_date" gorm:"type:date"`
	Location     string    `json:"location" gorm:"type:varchar(100)"`
	LocationType string    `json:"location_type" gorm:"type:enum('offline', 'online')"`
	Quota        int       `json:"quota" gorm:"type:int"`
	Price        float64   `json:"price" gorm:"type:decimal(8,2)"`
	Status       string    `json:"status" gorm:"type:enum('draft', 'published', 'cancelled')"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Sessions      []Session      `json:"session" gorm:"foreignKey:EventID"`
	Registrations []Registration `json:"registration" gorm:"foreignKey:EventID"`
}
