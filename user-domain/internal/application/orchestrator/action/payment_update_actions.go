package action

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-domain/internal/application/orchestrator/command"
	"user-domain/internal/application/outbound"
	"user-domain/internal/entity"

	"github.com/google/uuid"
)

// paymentUpdateExecution handles the execution of payment update
type paymentUpdateExecution struct {
	producer  outbound.Producer
	newUser   *entity.User
	commandID uuid.UUID
	isRan     bool
}

func (c *paymentUpdateExecution) Execute(ctx context.Context) error {
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
	c.isRan = true
	return nil
}

func (c *paymentUpdateExecution) Ran() bool {
	return c.isRan
}

func (c *paymentUpdateExecution) Name() string {
	return command.PAYMENT_UPDATE_EXECUTE
}

// NewPaymentUpdateExecution creates a new payment update execution command
func NewPaymentUpdateExecution(producer outbound.Producer, newUser *entity.User) command.Execution {
	return &paymentUpdateExecution{
		producer:  producer,
		newUser:   newUser,
		commandID: uuid.New(),
	}
}

// paymentUpdateCompensation handles the compensation (rollback) of payment update
type paymentUpdateCompensation struct {
	producer  outbound.Producer
	oldUser   *entity.User
	commandID uuid.UUID
	isRan     bool
}

func (c *paymentUpdateCompensation) Compensate(ctx context.Context) error {
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
	c.isRan = true
	return nil
}

func (c *paymentUpdateCompensation) Ran() bool {
	return c.isRan
}

func (c *paymentUpdateCompensation) Name() string {
	return command.PAYMENT_UPDATE_COMPENSATE
}

// NewPaymentUpdateCompensation creates a new payment update compensation command
func NewPaymentUpdateCompensation(producer outbound.Producer, oldUser *entity.User) command.Compensation {
	return &paymentUpdateCompensation{
		producer:  producer,
		oldUser:   oldUser,
		commandID: uuid.New(),
	}
}

// paymentUpdateVerification handles the verification of payment update
type paymentUpdateVerification struct {
	subscriber outbound.Subscriber
	isRan      bool
}

func (c *paymentUpdateVerification) Verify(ctx context.Context) error {
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
	c.isRan = true
	return nil
}

func (c *paymentUpdateVerification) Ran() bool {
	return c.isRan
}

func (c *paymentUpdateVerification) Name() string {
	return command.PAYMENT_UPDATE_VERIFICATION
}

// NewPaymentUpdateVerification creates a new payment update verification command
func NewPaymentUpdateVerification(subscriber outbound.Subscriber) command.Verification {
	return &paymentUpdateVerification{
		subscriber: subscriber,
	}
}

// paymentUpdateApproval handles the approval of payment update
type paymentUpdateApproval struct {
	producer  outbound.Producer
	userID    string
	commandID uuid.UUID
	isRan     bool
}

func (c *paymentUpdateApproval) Approve(ctx context.Context) error {
	bytes, err := json.Marshal(map[string]interface{}{
		"event":      "payment.approve",
		"user_id":    c.userID,
		"command_id": c.commandID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment approve event: %w", err)
	}

	if err := c.producer.Push(ctx, "payment.approve", bytes); err != nil {
		return fmt.Errorf("failed to publish payment approve: %w", err)
	}
	c.isRan = true
	return nil
}

func (c *paymentUpdateApproval) Ran() bool {
	return c.isRan
}

func (c *paymentUpdateApproval) Name() string {
	return command.PAYMENT_UPDATE_APPROVE
}

// NewPaymentUpdateApproval creates a new payment update approval command
func NewPaymentUpdateApproval(producer outbound.Producer, userID string) command.Approval {
	return &paymentUpdateApproval{
		producer:  producer,
		userID:    userID,
		commandID: uuid.New(),
	}
}
