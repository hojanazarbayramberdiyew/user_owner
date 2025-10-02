package repository

import (
	"context"
	"user_owner/internal/dto"
	"user_owner/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.CreateUserReq) error
	GetUsers(ctx context.Context) ([]model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// GetUsers implements UserRepository.
func (r *userRepository) GetUsers(ctx context.Context) ([]model.User, error) {

	query := "SELECT id,name,password,phone_number,email,logo,created_at,updated_at FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Password, &u.PhoneNumber, &u.Email, &u.Logo, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil

}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *dto.CreateUserReq) error {

	query := "INSERT INTO users(name,password,phone_number,email,logo) VALUES($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at"
	err := r.db.QueryRow(ctx, query,
		user.Name,
		user.Password,
		user.PhoneNumber,
		user.Email,
		user.Logo,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}
