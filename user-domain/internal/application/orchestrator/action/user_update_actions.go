package action

import (
	"context"
	"encoding/json"
	"fmt"
	"user-domain/internal/application/orchestrator/command"
	"user-domain/internal/application/outbound"
	"user-domain/internal/entity"

	"github.com/google/uuid"
)

// userUpdateApproval handles the approval of user update
type userUpdateApproval struct {
	producer  outbound.Producer
	userID    string
	commandID uuid.UUID
	isRan     bool
}

func (c *userUpdateApproval) Approve(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":      "user.approve",
		"user_id":    c.userID,
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal user approve event: %w", err)
	}
	if err := c.producer.Push(ctx, "user.approve", bytes); err != nil {
		return fmt.Errorf("failed to publish user approve: %w", err)
	}
	c.isRan = true
	return nil
}

func (c *userUpdateApproval) Ran() bool {
	return c.isRan
}

func (c *userUpdateApproval) Name() string {
	return command.USER_UPDATE_APPROVE
}

// NewUserUpdateApproval creates a new user update approval command
func NewUserUpdateApproval(producer outbound.Producer, userID string) command.Approval {
	return &userUpdateApproval{
		producer:  producer,
		userID:    userID,
		commandID: uuid.New(),
	}
}

// userUpdateCompensation handles the compensation (rollback) of user update
type userUpdateCompensation struct {
	producer  outbound.Producer
	oldUser   *entity.User
	commandID uuid.UUID
	isRan     bool
}

func (c *userUpdateCompensation) Compensate(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":      "user.rollback",
		"user":       c.oldUser,
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal user rollback event: %w", err)
	}

	if err := c.producer.Push(ctx, "user.rollback", bytes); err != nil {
		return fmt.Errorf("failed to publish user rollback: %w", err)
	}
	c.isRan = true
	return nil
}

func (c *userUpdateCompensation) Ran() bool {
	return c.isRan
}

func (c *userUpdateCompensation) Name() string {
	return command.USER_UPDATE_COMPENSATE
}

// NewUserUpdateCompensation creates a new user update compensation command
func NewUserUpdateCompensation(producer outbound.Producer, oldUser *entity.User) command.Compensation {
	return &userUpdateCompensation{
		producer:  producer,
		oldUser:   oldUser,
		commandID: uuid.New(),
	}
}
