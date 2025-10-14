package applicationinbound

import (
	"net/http"
)

type UserApi interface {
	PostUsers(w http.ResponseWriter, r *http.Request)
	PutUsersUserId(w http.ResponseWriter, r *http.Request, userID string)
	GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string)
	// List(offset int, limit int) (*dto.UsersResponse, error)
	DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId string)
}
