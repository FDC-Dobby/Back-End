package model

import (
	"github.com/google/uuid"
	"time"
)

type Loc struct {
	ID        uuid.UUID `json:"loc-id"`
	Name      string    `json:"name" validate:"required,min=1,max=20,excludesall=;"`
	Latitude  float64   `json:"latitude" validate:"required"`
	Longitude float64   `json:"longitude" validate:"required"`
	Runway    int64     `json:"runway"`
	Elevator  int64     `json:"elevator"`
	Parking   int64     `json:"parking"`
	Restroom  int64     `json:"restroom"`
	Block     int64     `json:"block"`
	Guide     int64     `json:"guide"`
	Review    []string  `json:"review"`
	CreatedAt time.Time `json:"created-at"`
}
