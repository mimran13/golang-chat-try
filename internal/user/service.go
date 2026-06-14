package user

import (
	"context"
	"log/slog"

	"go-chat/internal/apperror"
	"go-chat/internal/user/dto"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, input CreateUserInput) (User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type authService interface {
	GenerateToken(userID int64) (string, error)
}

type UserService struct {
	repo userRepository
	authService authService
}

func NewUserService(repo userRepository, auth authService) *UserService {
	return &UserService{repo: repo, authService: auth}
}

func (svc *UserService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		response[i] = toUserResponse(u)
	}

	return response, nil
}

func (svc *UserService) RegisterUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, apperror.Internal("failed to hash password")
	}

	input := CreateUserInput{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	created, err := svc.repo.Create(ctx, input)
	if err != nil {
		return dto.UserResponse{}, err
	}

	slog.Info("user created", "id", created.ID)

	return toUserResponse(created), nil
}

func (svc *UserService) LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	
	user, err := svc.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.LoginResponse{}, apperror.Internal("something happened")
	}

	if user == nil {
		return dto.LoginResponse{}, apperror.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
          return dto.LoginResponse{}, apperror.Unauthorized("invalid credentials")
      }

	token, tokenErr := svc.authService.GenerateToken(user.ID)
	if(tokenErr != nil) {
		return dto.LoginResponse{}, apperror.Internal("error generating token") 
	}

	return dto.LoginResponse{
		Token: token,
		User: toUserResponse(*user),
	}, nil
}

func toUserResponse(u User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}


