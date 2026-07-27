package room

import "time"

type RoomType string

const (
	RoomTypeDirect RoomType = "direct"
	RoomTypeGroup  RoomType = "group"
)

type Room struct {
	ID        int64
	Name      *string
	Type      RoomType
	CreatedBy int64
	CreatedAt time.Time
}

type RoomMember struct {
	RoomID   int64
	UserID   int64
	JoinedAt time.Time
}

type CreateRoomInput struct {
	Name      *string
	Type      RoomType
	CreatedBy int64
	MemberIDs []int64
}
