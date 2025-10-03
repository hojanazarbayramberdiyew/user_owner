package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Email       string    `json:"email" binding:"required"`
	Logo        string    `json:"logo"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Password    string    `json:"_"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Logo        string    `json:"logo"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	City   string  `json:"city" binding:"required"`
	Region string  `json:"region" binding:"required"`
	Price  float64 `json:"price" binding:"required,min=0"`
}

type OrderResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	UserID    string    `json:"user_id"`
	UserPhone string    `json:"user_phone"`
	City      string    `json:"city"`
	Region    string    `json:"region"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
