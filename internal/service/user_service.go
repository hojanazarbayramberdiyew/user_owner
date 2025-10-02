package service

import (
	"context"
	"user_owner/internal/dto"
	"user_owner/internal/model"
	"user_owner/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *dto.CreateUserReq) error
	GetUsers(ctx context.Context) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// GetUsers implements UserService.
func (s *userService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetUsers(ctx)
}

// CreateUser implements UserService.
func (s *userService) CreateUser(ctx context.Context, user *dto.CreateUserReq) error {
	return s.repo.CreateUser(ctx, user)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
