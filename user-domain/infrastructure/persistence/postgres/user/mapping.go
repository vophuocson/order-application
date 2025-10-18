package postgres

import (
	"user-domain/infrastructure/persistence/postgres/model"
	"user-domain/internal/entity"
)

func CreateRepoEntityFromUserEntity(e *entity.User) *model.User {
	return &model.User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
		Phone: e.Phone,
	}
}

func CreateUserEntityFromUserModel(e *model.User) *entity.User {
	return &entity.User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
		Phone: e.Phone,
	}
}

func CreateUsersEntityFromUsesrModel(usersModel []*model.User) []*entity.User {
	var result []*entity.User
	for _, u := range usersModel {
		result = append(result, CreateUserEntityFromUserModel(u))
	}
	return result
}
