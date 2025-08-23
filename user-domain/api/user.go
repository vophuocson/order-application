package api

import (
	"context"
	"user-domain/api/dto"
	"user-domain/api/inbound"
	"user-domain/internal/entity"
	"user-domain/internal/inport"
)

type user struct {
	sv inport.UserService
}

func (api *user) Create(userReq *dto.UserPost) error {
	userEntity := entity.User{}
	userReq.MapTo(&userEntity)
	ctx := context.Background()
	err := api.sv.CreateUser(ctx, &userEntity)
	return err
}

func (api *user) Update(userID string, userReq *dto.UserPut) error {
	userEntity := entity.User{ID: userID}
	userReq.MapTo(&userEntity)
	ctx := context.Background()
	err := api.sv.UpdateUser(ctx, &userEntity)
	return err
}

func (api *user) Get(userID string) (*dto.UserResponse, error) {
	ctx := context.Background()
	userEntity, err := api.sv.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := dto.UserResponse{}
	res.GetFrom(userEntity)
	return &res, nil
}

func (api *user) List(offset int, limit int) (*dto.UsersResponse, error) {
	ctx := context.Background()
	eUsers, err := api.sv.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	users := dto.UsersResponse{}
	users.GetFrom(eUsers)
	return &users, nil
}

func (api *user) Delete(userID string) error {
	ctx := context.Background()
	err := api.sv.DeleteUser(ctx, userID)
	return err
}

func NewApiUser(sv inport.UserService) inbound.UserApi {
	return &user{sv: sv}
}
