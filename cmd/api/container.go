package main

import (
	"database/sql"

	"go-chat/internal/handler"
	"go-chat/internal/user"
)

type Container struct {
	HealthHandler *handler.HealthHandler
	UserHandler *handler.UserHandler
}

func NewContainer(db *sql.DB) *Container {
	healthHandler := handler.NewHealthHandler(db)
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)

	return &Container{
		HealthHandler: healthHandler,
		UserHandler: handler.NewUserHandler(userSvc),
	}
}
