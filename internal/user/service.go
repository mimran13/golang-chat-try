package user

import (
	"context"
	"go-chat/internal/user/dto"
)

type userRepository interface {
	GetAll(ctx context.Context) ([]User, error)
}

type UserService struct {
 repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) GetAllUsers(ctx context.Context) ([]dto.GetUsersResponseDto, error) {
	users, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dto.GetUsersResponseDto, len(users))
	for i, u := range users {
		response[i] = dto.GetUsersResponseDto{
              ID:        u.ID,                                                                                                                                                                            
              Username:  u.Username,                                                                                                                                                                    
              Email:     u.Email,             
              CreatedAt: u.CreatedAt,     
          }
	}

	return response, nil;
}