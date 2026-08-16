package room

import (
	"context"
	"database/sql"
	"errors"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, input CreateRoomInput) (Room, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Room{}, err
	}

	defer func() { _ = tx.Rollback() }()

	result, err := tx.ExecContext(ctx,
		"INSERT INTO rooms (name, type, created_by) VALUES (?, ?, ?)", input.Name, input.Type, input.CreatedBy)
	if err != nil {
		return Room{}, err
	}

	roomID, err := result.LastInsertId()
	if err != nil {
		return Room{}, err
	}

	for _, memberID := range input.MemberIDs {
		_, err := tx.ExecContext(ctx, "INSERT INTO room_members (room_id, user_id) VALUES (?, ?)", roomID, memberID)
		if err != nil {
			return Room{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Room{}, err
	}

	var room Room
	err = r.db.QueryRowContext(ctx,
		"SELECT id, name, type, created_by, created_at FROM rooms WHERE id = ?",
		roomID,
	).Scan(&room.ID, &room.Name, &room.Type, &room.CreatedBy, &room.CreatedAt)
	if err != nil {
		return Room{}, err
	}

	return room, nil

}

func (r *RoomRepository) ListForUser(ctx context.Context, userID int64) ([]Room, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT r.id, r.name, r.type, r.created_by, r.created_at
        FROM rooms r
        INNER JOIN room_members rm ON rm.room_id = r.id
        WHERE rm.user_id = ?
        ORDER BY r.created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	rooms := make([]Room, 0)

	for rows.Next() {
		room := Room{}

		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Type,
			&room.CreatedBy,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil

}

func (r *RoomRepository) AddMember(ctx context.Context, roomID, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO room_members (room_id, user_id) VALUES (?, ?)",
		roomID, userID,
	)
	return err
}

func (r *RoomRepository) IsMember(ctx context.Context, roomID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = ? AND user_id = ?)",
		roomID, userID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *RoomRepository) FindByID(ctx context.Context, id int64) (*Room, error) {
	var room Room
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, type, created_by, created_at FROM rooms WHERE id = ?",
		id,
	).Scan(&room.ID, &room.Name, &room.Type, &room.CreatedBy, &room.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}
