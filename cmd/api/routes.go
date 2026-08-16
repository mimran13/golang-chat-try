package main

import (
	"sync/atomic"
	"time"

	"go-chat/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

const maxRequestBodyBytes = 1 << 20 // 1 MiB

func setupRoutes(c *Container, shuttingDown *atomic.Bool, allowedOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RequestIDMiddleware)
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middlewares.BodyLimitMiddleware(maxRequestBodyBytes))
	r.Use(middlewares.RateLimitMiddleware(100, 20))
	r.Use(middlewares.ShutdownMiddleware(shuttingDown))

	r.Get("/health", c.HealthHandler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", c.UserHandler.CreateUser)
		r.Post("/auth/login", c.UserHandler.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(c.AuthService))
			r.Get("/users", c.UserHandler.GetAll)
			r.Post("/rooms", c.RoomHandler.CreateRoom)
			r.Get("/rooms", c.RoomHandler.ListRooms)
			r.Post("/rooms/{id}/members", c.RoomHandler.AddMember)
		})
	})

	return r
}
