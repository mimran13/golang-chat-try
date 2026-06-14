package handler

import (
	"context"
	"errors"
	"net/http"

	"go-chat/internal/apperror"
	"go-chat/internal/response"
	"go-chat/internal/user/dto"
	"go-chat/internal/validation"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	RegisterUser(ctx context.Context, createUser dto.CreateUserRequest) (dto.UserResponse, error)
	LoginUser(ctx context.Context, login dto.LoginRequest) (dto.LoginResponse, error)
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAllUsers(r.Context())
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteAppError(w, appErr)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	response.WriteSuccess(w, http.StatusOK, users, nil)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := validation.ParseAndValidate(r, &req); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteAppError(w, appErr)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	created, err := h.svc.RegisterUser(r.Context(), req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteAppError(w, appErr)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, created, nil)
}


func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := validation.ParseAndValidate(r, &req); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteAppError(w, appErr)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	userLogin, err := h.svc.LoginUser(r.Context(), req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteAppError(w, appErr)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	response.WriteSuccess(w, http.StatusOK, userLogin, nil)
}
