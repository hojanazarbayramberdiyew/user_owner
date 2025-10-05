package service

import (
	"context"
	"fmt"
	"os"
	"user_owner/internal/dto"
	"user_owner/internal/model"
	"user_owner/internal/repository"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userPhone string, req *dto.CreateOrderRequest) (*model.Order, error)
	GenerateQRCodeFromCode(orderCode string) (string, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
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

	_, err := s.GenerateQRCodeFromCode(order.Code)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GenerateQRCodeFromCode(orderCode string) (string, error) {

	if err := os.MkdirAll("qrcodes", 0755); err != nil {
		return "", fmt.Errorf("directory döretmekde näsazlyk: %w", err)
	}
	filename := fmt.Sprintf("qrcodes/%s.png", orderCode)

	err := qrcode.WriteFile(orderCode, qrcode.Medium, 256, filename)
	if err != nil {
		return "", fmt.Errorf("QR code doretmekde nasazlyk: %w", err)
	}

	return filename, nil

}

func (s *orderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	return s.orderRepo.GetAllOrders(ctx)
}
