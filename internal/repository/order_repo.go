package repository

import (
	"context"
	"fmt"
	"user_owner/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func generateOrderCode() string {
	return uuid.New().String()[:5]
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	query := `
        INSERT INTO orders (code, user_id, user_phone, city, region, price)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

	order.Code = generateOrderCode()

	err := r.db.QueryRow(ctx, query,
		order.Code,
		order.UserID,
		order.UserPhone,
		order.City,
		order.Region,
		order.Price,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("order döretmekde näsazlyk: %w", err)
	}

	return nil
}
