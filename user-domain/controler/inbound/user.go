package inbound

import (
	"context"
	"user-domain/controler/handler"
)

type UserApi interface {
	PostUsers(ctx context.Context, request handler.PostUsersRequestObject) (handler.PostUsersResponseObject, error)
	// Update(userID string, userReq *dto.UserPut) error
	// Get(userID string) (*dto.UserResponse, error)
	// List(offset int, limit int) (*dto.UsersResponse, error)
	// Delete(userID string) error
}
