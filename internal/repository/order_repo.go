package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"user_owner/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetAllOrders(ctx context.Context) ([]model.Order, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func generateOrderCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 5

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
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

func (r *orderRepository) GetAllOrders(ctx context.Context) ([]model.Order, error) {

	query := "SELECT id,code,user_id,user_phone,city,region,price,created_at,updated_at FROM orders"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("order alyp bolmady: %w", err)
	}

	defer rows.Close()
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.Code, &order.UserID, &order.UserPhone, &order.City, &order.Region, &order.Price, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error on scan: %w", err)
		}
		orders = append(orders, order)

	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return orders, nil

}
