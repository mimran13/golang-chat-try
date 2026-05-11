package main

import (
	"database/sql"
	"go-chat/internal/handler"
	"go-chat/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	return r
}
