package orchestrator

import (
	"context"
	"user-domain/internal/application/orchestrator/workflow"
	"user-domain/internal/application/orchestrator/workflow/observer"
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
	// Generate workflow execution ID for tracing
	workflowID := "user_updation_" + uuid.New().String()
	o.logger.WithContext(ctx)

	o.logger.Info("Starting user update orchestration | workflow_id=%s | user_id=%s",
		workflowID, newUser.ID)

	// Create workflow with logging observer
	loggingObserver := observer.NewLoggingObserver(o.logger)
	updateWorkflow := workflow.NewUpdationUserWorkflowWithObservers(
		o.producer,
		o.subscriber,
		revertUser,
		newUser,
		loggingObserver,
	)

	// Create workflow context for distributed tracing
	workflowContext := &observer.WorkflowContext{
		WorkflowID:   workflowID,
		WorkflowType: "user_updation",
		EntityID:     newUser.ID,
	}

	// Inject context into workflow for tracing
	if wf, ok := updateWorkflow.(*workflow.UserUpdationWorkflow); ok {
		wf.SetContext(workflowContext)
	}

	// Execute workflow
	err := o.workflowRuner.Execute(ctx, workflowID, USER_UPDATION, updateWorkflow.Run)
	if err != nil {
		o.logger.Error("User update workflow failed | workflow_id=%s | user_id=%s | error=%v",
			workflowID, newUser.ID, err)
		return err
	}

	o.logger.Info("User update orchestration completed successfully | workflow_id=%s | user_id=%s",
		workflowID, newUser.ID)
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
