package middlewares

import (
	"context"
	"net/http"
	"strings"

	"go-chat/internal/response"
)

type tokenParser interface {
	ParseToken(token string) (int64, error)
}

type ctxKey string

const UserIDKey ctxKey = "user_id"

func AuthMiddleware(parser tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing token")
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			userID, err := parser.ParseToken(token)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}
