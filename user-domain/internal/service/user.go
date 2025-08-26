package service

import (
	"context"
	"user-domain/internal/entity"
	inport "user-domain/internal/inport"
	"user-domain/internal/outport"
)

type user struct {
	repo   outport.UserRepository
	logger outport.Logger
}

func (u *user) CreateUser(ctx context.Context, user *entity.User) error {
	// implement bussiness logic here
	err := u.repo.CreateUser(ctx, user)
	return err
}

func (u *user) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	// implement bussiness logic here
	userRes, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return userRes, err
}

func (u *user) UpdateUser(ctx context.Context, user *entity.User) error {
	// implement bussiness logic here
	err := u.repo.UpdateUser(ctx, user)
	return err
}

func (u *user) DeleteUser(ctx context.Context, id string) error {
	// implement bussiness logic here
	err := u.repo.DeleteUser(ctx, id)
	return err
}

func (u *user) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	// implement bussiness logic here
	entities, err := u.repo.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func NewUserService(r outport.UserRepository, logger outport.Logger) inport.UserService {
	return &user{repo: r, logger: logger}
}
