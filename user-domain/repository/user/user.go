package userrepo

import (
	"context"
	"fmt"
	"user-domain/internal/entity"
	"user-domain/repository"
	"user-domain/repository/dao"
	"user-domain/repository/model"

	"gorm.io/driver/postgres"
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

func NewUserRepo(cfg repository.DBConfig) (*UserRepo, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.GetHostName(), cfg.GetPort(), cfg.GetUser(), cfg.GetDBName(), cfg.GetPassword(), cfg.GetSSLMode())
	gormInstance, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	query := dao.Use(gormInstance)
	return &UserRepo{
		query: *query,
	}, nil
}
