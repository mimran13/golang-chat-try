package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-chat/internal/apperror"
)

type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message,omitempty"`
	Error   string            `json:"error,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message string, errorMessage string) {
	writeJSON(w, status, ErrorResponse{
		Status:  status,
		Message: message,
		Error:   errorMessage,
	})
}

func WriteAppError(w http.ResponseWriter, e *apperror.AppError) {
	writeJSON(w, e.Status, ErrorResponse{
		Status:  e.Status,
		Message: e.Code,
		Error:   e.Message,
		Fields:  e.Fields,
	})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, data)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

type SuccessResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, status int, data any, meta any) {
	WriteJSON(w, status, SuccessResponse{Data: data, Meta: meta})
}
