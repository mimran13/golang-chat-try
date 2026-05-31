package main

import (
	"go-chat/internal/middlewares"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func setupRoutes(c *Container) *chi.Mux {
	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.TimeoutMiddleware(5 * time.Second))
	r.Use(middlewares.RequestIDMiddleware)

	r.Get("/health", c.HealthHandler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", c.UserHandler.GetAll)
	})

	return r
}
