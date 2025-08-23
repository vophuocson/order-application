package inbound

import "user-domain/api/dto"

type UserApi interface {
	Create(userReq *dto.UserPost) error
	Update(userID string, userReq *dto.UserPut) error
	Get(userID string) (*dto.UserResponse, error)
	List(offset int, limit int) (*dto.UsersResponse, error)
	Delete(userID string) error
}
