package service

import (
	"context"
	"mime/multipart"
	"time"
	"user_owner/internal/config"
	"user_owner/internal/dto"
	"user_owner/internal/model"
	"user_owner/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, user *dto.CreateUserReq) error
	GetUsers(ctx context.Context) ([]model.User, error)
	Login(ctx context.Context, loginReq *dto.LoginReq) (*dto.LoginResponse, error)
	UpdateUserLogo(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error
}

type userService struct {
	repo        repository.UserRepository
	config      *config.Config
	fileService FileService
}

// GetUsers implements UserService.
func (s *userService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetUsers(ctx)
}

// CreateUser implements UserService.
func (s *userService) CreateUser(ctx context.Context, user *dto.CreateUserReq) error {
	user.ID = uuid.Nil
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
	return s.repo.CreateUser(ctx, user)
}

func NewUserService(repo repository.UserRepository, cfg *config.Config, fileService FileService) UserService {
	return &userService{repo: repo, config: cfg, fileService: fileService}
}

func (s *userService) Login(ctx context.Context, loginReq *dto.LoginReq) (*dto.LoginResponse, error) {

	user, err := s.repo.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  user,
	}, nil

}

func (s *userService) generateJWT(user *dto.UserResponse) (string, error) {

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Name,
		"email":      user.Email,
		"user_phone": user.PhoneNumber,
		"exp":        time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration)).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *userService) UpdateUserLogo(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error {

	logoPath, err := s.fileService.UploadUserLogo(ctx, userID, fileHeader)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserLogo(ctx, userID, logoPath)
	if err != nil {
		return err
	}

	return nil

}
