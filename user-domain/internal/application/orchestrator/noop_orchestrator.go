package orchestrator

import (
	"context"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

// noopOrchestrator is a no-op implementation of WorkflowOrchestrator
// Use this when you don't need workflow orchestration (e.g., for simple CRUD operations)
type noopOrchestrator struct{}

// ExecuteUserUpdation is a no-op implementation that does nothing
func (n *noopOrchestrator) ExecuteUserUpdation(ctx context.Context, revertUser, newUser *entity.User) error {
	// No-op: workflow orchestration is disabled
	return nil
}

// NewNoopOrchestrator creates a no-op orchestrator that does nothing
// This is useful for simple scenarios where workflow orchestration is not needed
func NewNoopOrchestrator() outport.WorkflowOrchestrator {
	return &noopOrchestrator{}
}

