package inbound

import (
	"net/http"
	"user-domain/internal/application/controller/parameter"
)

type UserApi interface {
	PostUsers(w http.ResponseWriter, r *http.Request)
	PutUsersUserId(w http.ResponseWriter, r *http.Request, userID string)
	GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string)
	DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId string)
	GetUsers(w http.ResponseWriter, r *http.Request, paramObj parameter.UserQueryParams)
}
