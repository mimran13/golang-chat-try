package shared

import (
	"encoding/json"
	"go-chat/internal/apperror"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

 var Validate = validator.New()

func ParseAndValidate(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		slog.Error("failed to decode request", "error", err)
		return apperror.BadRequest("Invalid request body")	
	}

	if err := Validate.Struct(dto); err != nil {
          return apperror.BadRequest(err.Error())
      }
      return nil
}