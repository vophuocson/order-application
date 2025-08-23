package router

import (
	"encoding/json"
	"net/http"
	"user-domain/api/dto"
	"user-domain/api/inbound"
)

type UserHandler struct {
	h inbound.UserApi
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto dto.UserPost
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return
	}
	defer r.Body.Close()
	err = h.h.Create(&userDto)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.h.List(0, 1)
	if err != nil {
		return
	}
	_ = users
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.h.Get("id")
	if err != nil {
		return
	}
	_ = user
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.h.Delete("id")
	if err != nil {
		return
	}
}
