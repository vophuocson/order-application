package userrepo

import (
	"context"
	"user-domain/internal/entity"
	"user-domain/repository/dao"
	"user-domain/repository/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	query dao.Query
}

func (d *UserRepo) CreateUser(ctx context.Context, user *entity.User) error {
	u := CreateRepoEntityFromUserEntity(user)
	userQery := d.query.User
	return userQery.WithContext(ctx).Create(u)
}

func (d *UserRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	u := CreateRepoEntityFromUserEntity(user)
	userQery := d.query.User
	_, err := userQery.WithContext(ctx).Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) DeleteUser(ctx context.Context, id string) error {
	userQery := d.query.User
	_, err := userQery.Delete(&model.User{ID: id})
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userQery := d.query.User
	userM, err := userQery.Where(userQery.ID.Eq(id)).First()
	if err != nil {
		return nil, nil
	}
	return CreateUserEntityFromUserModel(userM), nil
}

func (d *UserRepo) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	userQery := d.query.User
	usersModel, err := userQery.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, err
	}
	return CreateUsersEntityFromUsesrModel(usersModel), nil
}

func NewUserRepo(db *gorm.DB) (*UserRepo, error) {
	query := dao.Use(db)
	return &UserRepo{
		query: *query,
	}, nil
}
