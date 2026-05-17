package handler

import (
	"encoding/json"
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
      users, err := h.svc.GetAllUsers()
      if err != nil {
          http.Error(w, "failed to fetch users", http.StatusInternalServerError)
          return
      }

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(users)
  }