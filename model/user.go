package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `json:"user-id"`
	Username    string    `json:"username" validate:"required,min=1,max=20,excludesall=;"`
	Email       string    `json:"email" validate:"required,email"`
	PhoneNumber string    `json:"phone-number" validate:"required,e164"` //e164 format is used to save phone number
	Name        string    `json:"name" validate:"required,min=1,max=10,excludesall=;"`
	Password    string    `json:"-" validate:"required,min=8,max=30,excludesall=;"`
	Birthday    time.Time `json:"birthday" validate:"required"`
	CreatedAt   time.Time `json:"created-at"`
}
