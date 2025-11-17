package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-domain/internal/application/outbound"
	"user-domain/internal/entity"

	"github.com/google/uuid"
)

const (
	USER_UPDATE    = "UserUpdate"
	PAYMENT_UPDATE = "PaymentUpdate"
)

type UserUpdationActivity struct {
	commands []Command
}

func (uA UserUpdationActivity) GetCommands() []Command {
	return uA.commands
}

type CommandResponse struct {
	CommandName string
	Error       error
}

type userUpdateCommand struct {
	producer   outbound.Producer
	subscriber outbound.Subscriber
	oldUser    *entity.User
	newUser    *entity.User
	commandID  uuid.UUID
}

func (c *userUpdateCommand) Name() string {
	return USER_UPDATE
}

func (c *userUpdateCommand) ExecutePending(ctx context.Context) error { return nil }
func (c *userUpdateCommand) Verify(ctx context.Context) error         { return nil }
func (c *userUpdateCommand) Approve(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":      "user.approve",
		"user_id":    c.newUser.ID,
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal user approve event: %w", err)
	}
	if err := c.producer.Push(ctx, "user.approve", bytes); err != nil {
		return fmt.Errorf("failed to publish user approve: %w", err)
	}
	return nil
}

func (c *userUpdateCommand) Compensate(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":      "user.rollback",
		"user_id":    c.oldUser.ID,
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal user rollback event: %w", err)
	}

	if err := c.producer.Push(ctx, "user.rollback", bytes); err != nil {
		return fmt.Errorf("failed to publish user rollback: %w", err)
	}

	return nil
}

type paymentUpdateCommand struct {
	producer   outbound.Producer
	subscriber outbound.Subscriber
	oldUser    *entity.User
	newUser    *entity.User
	commandID  uuid.UUID
}

func (c *paymentUpdateCommand) Name() string {
	return PAYMENT_UPDATE
}

func (c *paymentUpdateCommand) ExecutePending(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":   "payment.pending",
		"user_id": c.newUser.ID,
		"data": map[string]string{
			"name":    c.newUser.Name,
			"email":   c.newUser.Email,
			"phone":   c.newUser.Phone,
			"address": c.newUser.Address,
		},
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment pending event: %w", err)
	}
	if err := c.producer.Push(ctx, "payment.pending", bytes); err != nil {
		return fmt.Errorf("failed to publish payment pending: %w", err)
	}
	return nil
}

func (c *paymentUpdateCommand) Verify(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	responseBytes, err := c.subscriber.Consume(ctx, "payment.pending.response")
	if err != nil {
		return fmt.Errorf("failed to receive payment verification: %w", err)
	}
	var response VerificationResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payment verification: %w", err)
	}
	if !response.Accepted {
		return fmt.Errorf("payment service rejected pending data: %s", response.Message)
	}
	return nil
}

func (c *paymentUpdateCommand) Approve(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":   "payment.approve",
		"user_id": c.newUser.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment approve event: %w", err)
	}

	if err := c.producer.Push(ctx, "payment.approve", bytes); err != nil {
		return fmt.Errorf("failed to publish payment approve: %w", err)
	}

	return nil
}

func (c *paymentUpdateCommand) Compensate(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":   "payment.rollback",
		"user_id": c.oldUser.ID,
		"data": map[string]string{
			"name":    c.oldUser.Name,
			"email":   c.oldUser.Email,
			"phone":   c.oldUser.Phone,
			"address": c.oldUser.Address,
		},
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment rollback event: %w", err)
	}

	if err := c.producer.Push(ctx, "payment.rollback", bytes); err != nil {
		return fmt.Errorf("failed to publish payment rollback: %w", err)
	}

	return nil
}

func NewUpdationUserActivity(producer outbound.Producer,
	subscriber outbound.Subscriber,
	newUser, oldUser *entity.User) *UserUpdationActivity {
	activity := &UserUpdationActivity{
		commands: make([]Command, 0),
	}
	activity.commands = append(activity.commands, &userUpdateCommand{
		producer:   producer,
		subscriber: subscriber,
		oldUser:    oldUser,
		newUser:    newUser,
		commandID:  uuid.New(),
	})
	activity.commands = append(activity.commands, &paymentUpdateCommand{
		producer:   producer,
		subscriber: subscriber,
		oldUser:    oldUser,
		newUser:    newUser,
		commandID:  uuid.New(),
	})
	return activity
}
