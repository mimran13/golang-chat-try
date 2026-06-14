package user

import "time"

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type CreateUserInput struct {
	Username     string
	Email        string
	PasswordHash string
}

type VerifyUserInput struct {
	Username string
	Password string
}