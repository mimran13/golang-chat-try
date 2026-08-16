package room

import (
	"context"
	"log/slog"
	"slices"

	"go-chat/internal/apperror"
	"go-chat/internal/room/dto"
)

type roomRepository interface {
	Create(ctx context.Context, input CreateRoomInput) (Room, error)
	ListForUser(ctx context.Context, userID int64) ([]Room, error)
	AddMember(ctx context.Context, roomID, userID int64) error
	IsMember(ctx context.Context, roomID, userID int64) (bool, error)
	FindByID(ctx context.Context, id int64) (*Room, error)
}

type RoomService struct {
	repo roomRepository
}

func NewRoomService(repo roomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (svc *RoomService) CreateRoom(ctx context.Context, callerID int64, req dto.CreateRoomRequest) (dto.RoomResponse, error) {
	roomType := RoomType(req.Type)

	memberIDs := ensureCallerIncluded(callerID, req.MemberIDs)

	switch roomType {
	case RoomTypeDirect:
		if req.Name != nil {
			return dto.RoomResponse{}, apperror.BadRequest("direct rooms cannot have a name")
		}
		if len(memberIDs) != 2 {
			return dto.RoomResponse{}, apperror.BadRequest("direct rooms must have exactly 2 members")
		}
	case RoomTypeGroup:
		if req.Name == nil {
			return dto.RoomResponse{}, apperror.BadRequest("group rooms require a name")
		}
	default:
		return dto.RoomResponse{}, apperror.BadRequest("invalid room type")
	}

	input := CreateRoomInput{
		Name:      req.Name,
		Type:      roomType,
		CreatedBy: callerID,
		MemberIDs: memberIDs,
	}

	created, err := svc.repo.Create(ctx, input)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	slog.Info("room created", "id", created.ID, "type", created.Type, "created_by", callerID)

	return toRoomResponse(created), nil
}

func (svc *RoomService) ListRooms(ctx context.Context, callerID int64) ([]dto.RoomResponse, error) {
	rooms, err := svc.repo.ListForUser(ctx, callerID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.RoomResponse, len(rooms))
	for i, r := range rooms {
		response[i] = toRoomResponse(r)
	}
	return response, nil
}

func (svc *RoomService) AddMember(ctx context.Context, callerID, roomID, targetUserID int64) error {
	room, err := svc.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return apperror.NotFound("room not found")
	}

	isCallerMember, err := svc.repo.IsMember(ctx, roomID, callerID)
	if err != nil {
		return err
	}
	if !isCallerMember {
		return apperror.Forbidden("only members can add others to a room")
	}

	if room.Type == RoomTypeDirect {
		return apperror.BadRequest("cannot add members to a direct room")
	}

	alreadyMember, err := svc.repo.IsMember(ctx, roomID, targetUserID)
	if err != nil {
		return err
	}
	if alreadyMember {
		return apperror.Conflict("user is already a member")
	}

	if err := svc.repo.AddMember(ctx, roomID, targetUserID); err != nil {
		return err
	}

	slog.Info("member added", "room_id", roomID, "user_id", targetUserID, "added_by", callerID)
	return nil
}

func ensureCallerIncluded(callerID int64, ids []int64) []int64 {
	if slices.Contains(ids, callerID) {
		return ids
	}
	return append(ids, callerID)
}

func toRoomResponse(r Room) dto.RoomResponse {
	return dto.RoomResponse{
		ID:        r.ID,
		Name:      r.Name,
		Type:      string(r.Type),
		CreatedBy: r.CreatedBy,
		CreatedAt: r.CreatedAt,
	}
}
