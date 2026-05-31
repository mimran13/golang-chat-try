package middlewares

import (
	"go-chat/internal/response"
	"net/http"
	"sync/atomic"
)

func ShutdownMiddleware(shuttingDown *atomic.Bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shuttingDown.Load() {
				response.WriteError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "server is shutting down")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
