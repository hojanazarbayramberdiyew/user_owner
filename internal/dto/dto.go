package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Logo        string    `json:"logo"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
