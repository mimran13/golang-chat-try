package main

import (
	"database/sql"

	"go-chat/internal/handler"
	"go-chat/internal/user"
)

type Container struct {
	UserHandler *handler.UserHandler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)

	return &Container{
		UserHandler: handler.NewUserHandler(userSvc),
	}
}
