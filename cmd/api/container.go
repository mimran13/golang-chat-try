package main

import (
	"database/sql"

	"go-chat/internal/auth"
	"go-chat/internal/config"
	"go-chat/internal/handler"
	"go-chat/internal/user"
)

type Container struct {
	HealthHandler *handler.HealthHandler
	UserHandler   *handler.UserHandler
	AuthService   *auth.AuthService
}

func NewContainer(db *sql.DB, config *config.Config) *Container {
	healthHandler := handler.NewHealthHandler(db)
	authService := auth.NewAuthService(config.Auth)
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo, authService)

	return &Container{
		HealthHandler: healthHandler,
		UserHandler:   handler.NewUserHandler(userSvc),
		AuthService:   authService,
	}
}
