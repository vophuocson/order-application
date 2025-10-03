package userrepository

import (
	"context"
	userrepooutbound "user-domain/internal/application/repository/user/outbound"
	"user-domain/internal/entity"
	"user-domain/internal/outport"
)

type userRepo struct {
	userOutbound userrepooutbound.UserRepo
}

func (u *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	return u.userOutbound.CreateUser(ctx, user)
}

func (u *userRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.userOutbound.UpdateUser(ctx, user)
}

func (u *userRepo) DeleteUser(ctx context.Context, id string) error {
	return u.userOutbound.DeleteUser(ctx, id)
}

func (u *userRepo) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userOutbound.GetUserByID(ctx, id)
}

func (u *userRepo) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	return u.userOutbound.ListUsers(ctx, offset, limit)
}

func NewUserRepo(userOutbound userrepooutbound.UserRepo) outport.UserRepository {
	return &userRepo{
		userOutbound: userOutbound,
	}
}
