package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	UserID    uuid.UUID `json:"user_id"`
	UserPhone string    `json:"user_phone"`
	City      string    `json:"city"`
	Region    string    `json:"region"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}




