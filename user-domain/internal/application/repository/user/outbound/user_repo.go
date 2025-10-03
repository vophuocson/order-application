package userrepooutbound

import (
	"context"
	"user-domain/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)
}
