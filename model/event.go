package model

import (
	"time"

	"github.com/google/uuid"
)

type EventParam struct {
	EventID uuid.UUID `json:"-"`
}

type RegisterEvent struct {
	EventID uuid.UUID `json:"event_id"`
}

type CreateEventParam struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Location     string    `json:"location"`
	LocationType string    `json:"location_type"`
	Quota        int       `json:"quota"`
	Price        float64   `json:"price"`
}

type CreateEventResponse struct {
	EventID      uuid.UUID `json:"event_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Location     string    `json:"location"`
	LocationType string    `json:"location_type"`
	Quota        int       `json:"quota"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
}
