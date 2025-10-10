package model

import "github.com/google/uuid"

type RegistrationParam struct {
	RegistrationID uuid.UUID `json:"-"`
}
