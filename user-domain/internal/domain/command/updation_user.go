package command

import (
	"context"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type UpdationUser struct {
	repo       outport.UserRepository
	currenUser *entity.User
	oldUser    *entity.User
}

func (uU *UpdationUser) Execute(ctx context.Context) error {
	old, err := uU.repo.GetUserByID(ctx, uU.currenUser.ID)
	if err != nil {
		return err
	}
	uU.oldUser = old
	err = uU.repo.UpdateUser(ctx, uU.currenUser)
	if err != nil {
		return err
	}
	return nil
}

func (uU *UpdationUser) Undo(ctx context.Context) error {
	err := uU.repo.UpdateUser(ctx, uU.oldUser)
	if err != nil {
		return err
	}
	return nil
}

func NewUpdationUserCommand(repo outport.UserRepository, updatingUser *entity.User) Command {
	return &UpdationUser{
		repo:       repo,
		currenUser: updatingUser,
	}
}
