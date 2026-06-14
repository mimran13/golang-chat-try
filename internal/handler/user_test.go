package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-chat/internal/apperror"
	"go-chat/internal/user/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	getAllFunc   func(ctx context.Context) ([]dto.UserResponse, error)
	registerFunc func(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	loginFunc    func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	return m.getAllFunc(ctx)
}
func (m *mockUserService) RegisterUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	return m.registerFunc(ctx, req)
}
func (m *mockUserService) LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	return m.loginFunc(ctx, req)
}

func TestCreateUser_Success(t *testing.T) {
	svc := &mockUserService{
		registerFunc: func(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
			return dto.UserResponse{ID: 1, Username: req.Username, Email: req.Email}, nil
		},
	}
	h := NewUserHandler(svc)

	body, _ := json.Marshal(dto.CreateUserRequest{Username: "john", Email: "john@test.com", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Data dto.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, int64(1), resp.Data.ID)
	assert.Equal(t, "john", resp.Data.Username)
}

func TestCreateUser_InvalidBody(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{"username":"x"}`)))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUser_Success(t *testing.T) {
	svc := &mockUserService{
		loginFunc: func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
			return dto.LoginResponse{
				Token: "fake-jwt",
				User:  dto.UserResponse{ID: 1, Username: req.Username},
			}, nil
		},
	}
	h := NewUserHandler(svc)

	body, _ := json.Marshal(dto.LoginRequest{Username: "john", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.LoginUser(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data dto.LoginResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "fake-jwt", resp.Data.Token)
	assert.Equal(t, "john", resp.Data.User.Username)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	svc := &mockUserService{
		loginFunc: func(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
			return dto.LoginResponse{}, apperror.Unauthorized("invalid credentials")
		},
	}
	h := NewUserHandler(svc)

	body, _ := json.Marshal(dto.LoginRequest{Username: "john", Password: "wrongpassword"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.LoginUser(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
