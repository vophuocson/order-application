package userpersistence

import (
	"context"
	"user-domain/infrastructure/persistence/postgres/dao"
	"user-domain/infrastructure/persistence/postgres/model"
	applicationoutbound "user-domain/internal/application/outbound"
	"user-domain/internal/entity"

	"gorm.io/gorm"
)

type userRepo struct {
	query dao.Query
}

func (d *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	u := CreateRepoEntityFromUserEntity(user)
	userQery := d.query.User
	return userQery.WithContext(ctx).Create(u)
}

func (d *userRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	u := CreateRepoEntityFromUserEntity(user)
	userQery := d.query.User
	_, err := userQery.WithContext(ctx).Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (d *userRepo) DeleteUser(ctx context.Context, id string) error {
	userQery := d.query.User
	_, err := userQery.Delete(&model.User{ID: id})
	if err != nil {
		return err
	}
	return nil
}

func (d *userRepo) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userQery := d.query.User
	userM, err := userQery.Where(userQery.ID.Eq(id)).First()
	if err != nil {
		return nil, nil
	}
	return CreateUserEntityFromUserModel(userM), nil
}

func (d *userRepo) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	userQery := d.query.User
	usersModel, err := userQery.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, err
	}
	return CreateUsersEntityFromUsesrModel(usersModel), nil
}

func NewUserRepo(db *gorm.DB) applicationoutbound.UserRepo {
	query := dao.Use(db)
	return &userRepo{
		query: *query,
	}
}
