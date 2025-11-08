package command

import (
	"context"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type UpdationUserPayment struct {
	producer       outport.Producer
	compensateUser *entity.User
	commandUser    *entity.User
}

func (uP *UpdationUserPayment) Execute(ctx context.Context) error {
	return uP.producer.Push(ctx, "channel name", uP.commandUser)
}

func (uP *UpdationUserPayment) Undo(ctx context.Context) error {
	return uP.producer.Push(ctx, "channel name", uP.compensateUser)
}

func NewUpdationUserPayment(producer outport.Producer, commandUser *entity.User, compensateUser *entity.User) Command {
	return &UpdationUserPayment{
		producer:       producer,
		commandUser:    commandUser,
		compensateUser: compensateUser,
	}
}
