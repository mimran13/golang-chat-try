package room

import (
	"context"
	"database/sql"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) * RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) createRoom(ctx context.Context, input CreateRoomInput ) (Room, error) {
 result, err := r.db.ExecContext(ctx, 
	"INSERT INTO room (name, type, created_by) VALUES (?, ?, ?)", input.Name, input.Type, input.CreatedBy)
	if(err != nil) {
		return Room{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Room{}, err
	}

	var u Room
	err = r.db.QueryRowContext(ctx,
		"SELECT id, name, type, created_by, created_at FROM room WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.Type, &u.CreatedBy, &u.CreatedAt)
	if err != nil {
		return Room{}, err
	}

	return u, nil

}