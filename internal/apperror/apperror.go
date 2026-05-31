package apperror

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(message string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: message, Status: http.StatusNotFound}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: message, Status: http.StatusUnauthorized}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: "FORBIDDEN", Message: message, Status: http.StatusForbidden}
}

func BadRequest(message string) *AppError {
	return &AppError{Code: "BAD_REQUEST", Message: message, Status: http.StatusBadRequest}
}

func Internal(message string) *AppError {
	return &AppError{Code: "INTERNAL_SERVER_ERROR", Message: message, Status: http.StatusInternalServerError}
}

func Conflict(message string) *AppError {
	return &AppError{Code: "CONFLICT", Message: message, Status: http.StatusConflict}
}
