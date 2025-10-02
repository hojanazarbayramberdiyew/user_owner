package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Password    string
	PhoneNumber string
	Email       string
	Logo        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
