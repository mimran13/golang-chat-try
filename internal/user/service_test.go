package user

import (
	"context"
	"errors"
	"testing"

	"go-chat/internal/apperror"
	"go-chat/internal/user/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	getAllFunc         func(ctx context.Context) ([]User, error)
	createFunc         func(ctx context.Context, input CreateUserInput) (User, error)
	findByUsernameFunc func(ctx context.Context, username string) (*User, error)
}

func (m *mockRepo) GetAll(ctx context.Context) ([]User, error) {
	return m.getAllFunc(ctx)
}
func (m *mockRepo) Create(ctx context.Context, input CreateUserInput) (User, error) {
	return m.createFunc(ctx, input)
}
func (m *mockRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	return m.findByUsernameFunc(ctx, username)
}

type mockAuth struct {
	generateFunc func(userID int64) (string, error)
}

func (m *mockAuth) GenerateToken(userID int64) (string, error) {
	return m.generateFunc(userID)
}

func TestRegisterUser_Success(t *testing.T) {
	repo := &mockRepo{
		createFunc: func(ctx context.Context, input CreateUserInput) (User, error) {
			assert.Equal(t, "john", input.Username)
			assert.Equal(t, "john@test.com", input.Email)
			assert.NotEqual(t, "password123", input.PasswordHash, "password must be hashed")
			return User{ID: 1, Username: input.Username, Email: input.Email}, nil
		},
	}
	svc := NewUserService(repo, &mockAuth{})

	res, err := svc.RegisterUser(context.Background(), dto.CreateUserRequest{
		Username: "john", Email: "john@test.com", Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "john", res.Username)
}

func TestRegisterUser_RepoError(t *testing.T) {
	repo := &mockRepo{
		createFunc: func(ctx context.Context, input CreateUserInput) (User, error) {
			return User{}, errors.New("db down")
		},
	}
	svc := NewUserService(repo, &mockAuth{})

	_, err := svc.RegisterUser(context.Background(), dto.CreateUserRequest{
		Username: "john", Email: "john@test.com", Password: "password123",
	})

	require.Error(t, err)
}

func TestLoginUser_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*User, error) {
			return &User{ID: 1, Username: "john", PasswordHash: string(hash)}, nil
		},
	}
	auth := &mockAuth{
		generateFunc: func(userID int64) (string, error) {
			assert.Equal(t, int64(1), userID)
			return "fake-jwt-token", nil
		},
	}
	svc := NewUserService(repo, auth)

	res, err := svc.LoginUser(context.Background(), dto.LoginRequest{
		Username: "john", Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "fake-jwt-token", res.Token)
	assert.Equal(t, int64(1), res.User.ID)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	repo := &mockRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*User, error) {
			return nil, nil
		},
	}
	svc := NewUserService(repo, &mockAuth{})

	_, err := svc.LoginUser(context.Background(), dto.LoginRequest{Username: "ghost", Password: "x"})

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	repo := &mockRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*User, error) {
			return &User{ID: 1, Username: "john", PasswordHash: string(hash)}, nil
		},
	}
	svc := NewUserService(repo, &mockAuth{})

	_, err := svc.LoginUser(context.Background(), dto.LoginRequest{Username: "john", Password: "wrong-password"})

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
}

func TestLoginUser_TokenGenerationFails(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*User, error) {
			return &User{ID: 1, Username: "john", PasswordHash: string(hash)}, nil
		},
	}
	auth := &mockAuth{
		generateFunc: func(userID int64) (string, error) {
			return "", errors.New("signing failed")
		},
	}
	svc := NewUserService(repo, auth)

	_, err := svc.LoginUser(context.Background(), dto.LoginRequest{Username: "john", Password: "password123"})

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", appErr.Code)
}
