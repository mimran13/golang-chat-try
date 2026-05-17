package main

import (
	"go-chat/internal/handler"
	"go-chat/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRoutes(c *Container) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)
	r.Get("/users", c.UserHandler.GetAll)

	return r
}
