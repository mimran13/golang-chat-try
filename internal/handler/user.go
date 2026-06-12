package handler

import (
	"context"
	"errors"
	"go-chat/internal/apperror"
	"go-chat/internal/response"
	"go-chat/internal/user/dto"
	"net/http"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
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
