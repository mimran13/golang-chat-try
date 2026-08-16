package main

import (
	"database/sql"

	"go-chat/internal/auth"
	"go-chat/internal/config"
	"go-chat/internal/handler"
	"go-chat/internal/room"
	"go-chat/internal/user"
)

type Container struct {
	HealthHandler *handler.HealthHandler
	UserHandler   *handler.UserHandler
	RoomHandler   *handler.RoomHandler
	AuthService   *auth.AuthService
}

func NewContainer(db *sql.DB, config *config.Config) *Container {
	healthHandler := handler.NewHealthHandler(db)
	authService := auth.NewAuthService(config.Auth)
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo, authService)
	roomRepo := room.NewRoomRepository(db)
	roomSvc := room.NewRoomService(roomRepo)

	return &Container{
		HealthHandler: healthHandler,
		UserHandler:   handler.NewUserHandler(userSvc),
		RoomHandler:   handler.NewRoomHandler(roomSvc),
		AuthService:   authService,
	}
}
