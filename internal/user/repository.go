package user

import (
	"context"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, username, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, input CreateUserInput) (User, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		input.Username, input.Email, input.PasswordHash,
	)
	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	var u User
	err = r.db.QueryRowContext(ctx,
		"SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, password_hash, created_at FROM users where username = ?", 
		username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if(errors.Is(err, sql.ErrNoRows)) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
