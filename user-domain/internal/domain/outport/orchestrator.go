package outport

import (
	"context"
	"user-domain/internal/entity"
)

type WorkflowOrchestrator interface {
	ExecuteUserUpdation(ctx context.Context, revertUser, newUser *entity.User) error
}
