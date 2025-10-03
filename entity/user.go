package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	RoleID    int       `json:"role_id"`
	GoogleID  *string   `json:"google_id" gorm:"uniqueIndex"`
	Picture   *string   `json:"picture"`
	Name      string    `json:"name" gorm:"type:varchar(100)"`
	Email     string    `json:"email" gorm:"type:varchar(100)"`
	Password  *string   `json:"password" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Events        []Event        `json:"events" gorm:"foreignKey:UserID"`
	Registrations []Registration `json:"registrations" gorm:"foreignKey:UserID"`
}
