package validation

import (
	"encoding/json"
	"errors"
	"go-chat/internal/apperror"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ParseAndValidate(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		slog.Debug("failed to decode request", "error", err)
		return apperror.BadRequest("invalid request body")
	}

	if err := validate.Struct(dto); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			for _, fe := range ve {
				fields[fe.Field()] = fe.Tag()
			}
			return apperror.BadRequestFields(fields)
		}
		return apperror.BadRequest("invalid request body")
	}

	return nil
}
