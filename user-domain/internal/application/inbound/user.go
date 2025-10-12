package applicationinbound

import (
	"net/http"
)

type UserApi interface {
	PostUsers(w http.ResponseWriter, r *http.Request)
	// Update(userID string, userReq *dto.UserPut) error
	GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string)
	// List(offset int, limit int) (*dto.UsersResponse, error)
	// Delete(userID string) error
}
