package dto

import (
	"user-domain/internal/entity"
)

type UserPost struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

func (u UserPost) MapTo(e *entity.User) {
	e.Name = u.Name
	e.Email = u.Email
	e.Phone = u.Phone
	e.Address = u.Address
}

type UserPut struct {
	Name    *string `json:"name,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

func (u UserPut) MapTo(e *entity.User) {
	if u.Name != nil {
		e.Name = *u.Name
	}
	if u.Phone != nil {
		e.Phone = *u.Phone
	}
	if u.Address != nil {
		e.Address = *u.Address
	}
}

type UserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

func (u *UserResponse) GetFrom(e *entity.User) {
	u.ID = e.ID
	u.Name = e.Name
	u.Email = e.Email
	u.Phone = e.Phone
	u.Address = e.Address
}

type UsersResponse struct {
	Item []*UserResponse `json:"item"`
}

func (u *UsersResponse) GetFrom(users []*entity.User) {
	for _, us := range users {
		uRes := UserResponse{}
		uRes.GetFrom(us)
		u.Item = append(u.Item, &uRes)
	}
}
