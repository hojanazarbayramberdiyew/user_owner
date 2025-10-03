package service

import (
	"context"
	"user_owner/internal/dto"
	"user_owner/internal/model"
	"user_owner/internal/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userPhone string, req *dto.CreateOrderRequest) (*model.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, userPhone string, req *dto.CreateOrderRequest) (*model.Order, error) {
	order := &model.Order{
		UserID:    userID,
		UserPhone: userPhone,
		City:      req.City,
		Region:    req.Region,
		Price:     req.Price,
	}

	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
