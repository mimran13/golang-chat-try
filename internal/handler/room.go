package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"go-chat/internal/apperror"
	"go-chat/internal/middlewares"
	"go-chat/internal/response"
	"go-chat/internal/room/dto"
	"go-chat/internal/validation"

	"github.com/go-chi/chi/v5"
)

type RoomService interface {
	CreateRoom(ctx context.Context, callerID int64, req dto.CreateRoomRequest) (dto.RoomResponse, error)
	ListRooms(ctx context.Context, callerID int64) ([]dto.RoomResponse, error)
	AddMember(ctx context.Context, callerID, roomID, targetUserID int64) error
}

type RoomHandler struct {
	svc RoomService
}

func NewRoomHandler(svc RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		response.WriteAppError(w, apperror.Unauthorized("missing user context"))
		return
	}

	var req dto.CreateRoomRequest
	if err := validation.ParseAndValidate(r, &req); err != nil {
		writeErr(w, err)
		return
	}

	created, err := h.svc.CreateRoom(r.Context(), callerID, req)
	if err != nil {
		writeErr(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, created, nil)
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		response.WriteAppError(w, apperror.Unauthorized("missing user context"))
		return
	}

	rooms, err := h.svc.ListRooms(r.Context(), callerID)
	if err != nil {
		writeErr(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, rooms, nil)
}

func (h *RoomHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		response.WriteAppError(w, apperror.Unauthorized("missing user context"))
		return
	}

	roomID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || roomID <= 0 {
		response.WriteAppError(w, apperror.BadRequest("invalid room id"))
		return
	}

	var req dto.AddMemberRequest
	if err := validation.ParseAndValidate(r, &req); err != nil {
		writeErr(w, err)
		return
	}

	if err := h.svc.AddMember(r.Context(), callerID, roomID, req.UserID); err != nil {
		writeErr(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusNoContent, nil, nil)
}

func writeErr(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.WriteAppError(w, appErr)
		return
	}
	response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong")
}
