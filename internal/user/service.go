package user

import (
	"context"
	"log/slog"

	"go-chat/internal/user/dto"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	RegisterUser(ctx context.Context, user RegisterUser) (User, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		response[i] = dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	}

	return response, nil
}

func (svc *UserService) RegisterUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, err
	}

	repoInput := RegisterUser{
		Username: req.Username,
		Email:    req.Email,
		Password: string(passwordHash),
	}

	created, err := svc.repo.RegisterUser(ctx, repoInput)
	if err != nil {
		return dto.UserResponse{}, err
	}

	slog.Info("user created", "id", created.ID)

	return dto.UserResponse{
		ID:        created.ID,
		Username:  created.Username,
		Email:     created.Email,
		CreatedAt: created.CreatedAt,
	}, nil
}
