package user

import (
	"context"
	"user-domain/internal/domain/inport"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type user struct {
	repo         outport.UserRepository
	logger       outport.Logger
	orchestrator outport.WorkflowOrchestrator
}

func (u *user) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userRes, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return userRes, nil
}

func (u *user) UpdateUser(ctx context.Context, user *entity.User) error {
	old, err := u.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	err = u.repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	err = u.orchestrator.ExecuteUserUpdation(ctx, user, old)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, id)
}

func (u *user) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	entities, err := u.repo.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func NewUserService(r outport.UserRepository, logger outport.Logger, orchestrator outport.WorkflowOrchestrator) inport.UserService {
	return &user{repo: r, logger: logger, orchestrator: orchestrator}
}
