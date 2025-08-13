package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID `json:"session_id" gorm:"type:varchar(36);primaryKey"`
	EventID   uuid.UUID `json:"event_id"`
	Title     string    `json:"title" gorm:"type:varchar(70)"`
	StartTime time.Time `json:"start_time" gorm:"type:datetime"`
	EndTime   time.Time `json:"end_time" gorm:"type:datetime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
