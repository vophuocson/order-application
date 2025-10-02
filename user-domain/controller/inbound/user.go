package inbound

import "net/http"

type UserApi interface {
	PostUsers(w http.ResponseWriter, r *http.Request)
	// Update(userID string, userReq *dto.UserPut) error
	// Get(userID string) (*dto.UserResponse, error)
	// List(offset int, limit int) (*dto.UsersResponse, error)
	// Delete(userID string) error
}
