package orchestrator

import (
	"context"
	"user-domain/internal/application/orchestrator/workflow"
	"user-domain/internal/application/outbound"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"

	"github.com/google/uuid"
)

type orchestrator struct {
	producer      outbound.Producer
	subscriber    outbound.Subscriber
	logger        outbound.Logger
	workflowRuner outbound.WorkflowRuner
}

func (o *orchestrator) ExecuteUserUpdation(ctx context.Context, revertUser, newUser *entity.User) error {
	o.logger.Info("Starting user update orchestration for user ID: %s", newUser.ID)

	// Create workflow with logging observer
	loggingObserver := workflow.NewLoggingObserver(o.logger)
	updateWorkflow := workflow.NewUpdationUserWorkflowWithObservers(
		o.producer,
		o.subscriber,
		revertUser,
		newUser,
		loggingObserver,
	)

	workflowID := "user_updation" + uuid.New().String()
	err := o.workflowRuner.Execute(ctx, workflowID, USER_UPDATION, updateWorkflow.Run)
	if err != nil {
		o.logger.Error("User update workflow failed: %v", err)
		return err
	}
	o.logger.Info("User update orchestration completed successfully for user ID: %s", newUser.ID)
	return nil
}

func NewWorkflowStarter(
	producer outbound.Producer,
	subscriber outbound.Subscriber,
	logger outbound.Logger,
	workflowRuner outbound.WorkflowRuner,
) outport.WorkflowOrchestrator {
	return &orchestrator{
		producer:      producer,
		subscriber:    subscriber,
		logger:        logger,
		workflowRuner: workflowRuner,
	}
}
