package handler

import (
	"errors"
	"go-chat/internal/apperror"
	"go-chat/internal/response"
	"go-chat/internal/user"
	"net/http"
)

type UserHandler struct {
	svc *user.UserService
}

func NewUserHandler(svc *user.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAllUsers(r.Context())
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.WriteError(w, appErr.Status, appErr.Code, appErr.Message)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
		return
	}

	response.WriteSuccess(w, http.StatusOK, users, nil)
}