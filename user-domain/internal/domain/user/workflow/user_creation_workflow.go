package workflow

import (
	"user-domain/internal/domain/inport"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type UserCreationWorkflow struct {
	repo     outport.UserRepository
	producer outport.Producer
	consumer inport.Consumer
	oldData  *entity.User
}
