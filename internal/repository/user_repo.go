package repository

import (
	"context"
	"fmt"
	"user_owner/internal/dto"
	"user_owner/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.CreateUserReq) error
	GetUsers(ctx context.Context) ([]model.User, error)
	FindByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	Login(ctx context.Context, loginReq *dto.LoginReq) (*dto.UserResponse, error)
	UpdateUserLogo(ctx context.Context, userID uuid.UUID, logoPath string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil

}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (r *userRepository) Login(ctx context.Context, loginReq *dto.LoginReq) (*dto.UserResponse, error) {
	user, err := r.FindByUsername(ctx, loginReq.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if err := CheckPassword(user.Password, loginReq.Password); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return user, nil
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

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("password hashlemekde nasazlyk: %v", err)
	}

	query := "INSERT INTO users(name,password,phone_number,email,logo) VALUES($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at"
	err = r.db.QueryRow(ctx, query,
		user.Name,
		hashedPassword,
		user.PhoneNumber,
		user.Email,
		user.Logo,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("user döretmekde näsazlyk: %w", err)
	}

	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	var user dto.UserResponse
	query := "SELECT id,name,password,phone_number,email,logo,created_at,updated_at FROM users WHERE name=$1"

	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Name, &user.Password, &user.PhoneNumber, &user.Email, &user.Logo, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found by name: %s", username)
		}
		return nil, fmt.Errorf("user alyp bolmady: %v", err)
	}

	return &user, nil

}

func (r *userRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE name = $1"

	err := r.db.QueryRow(ctx, query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("username barlap bolmady: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) UpdateUserLogo(ctx context.Context, userID uuid.UUID, logoPath string) error {

	query := "UPDATE users SET logo=$1,updated_at=Now() WHERE id=$2"

	_, err := r.db.Exec(ctx, query, logoPath, userID)
	if err != nil {
		return fmt.Errorf("user logosy tazelenmedi: %v", err)
	}
	return nil

}
