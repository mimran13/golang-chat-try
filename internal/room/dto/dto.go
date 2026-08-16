package dto

import "time"

type CreateRoomRequest struct {
	Type      string  `json:"type" validate:"required,oneof=direct group"`
	Name      *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	MemberIDs []int64 `json:"member_ids" validate:"required,min=1,dive,gt=0"`
}

type RoomResponse struct {
 	ID        int64     `json:"id"`
    Name      *string   `json:"name"`
    Type      string    `json:"type"`
    CreatedBy int64     `json:"created_by"`
    CreatedAt time.Time `json:"created_at"`
}

type AddMemberRequest struct {
    UserID int64 `json:"user_id" validate:"required,gt=0"`
}
