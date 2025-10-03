package userrepository

import (
	"context"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type userRepo struct {
	userOutbound applicationoutbound.UserRepo
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

func NewUserRepo(userOutbound applicationoutbound.UserRepo) domainoutport.UserRepository {
	return &userRepo{
		userOutbound: userOutbound,
	}
}
