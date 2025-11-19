package orchestrator

import (
	"context"
	"fmt"
	"user-domain/internal/application/orchestrator/workflow"
	"user-domain/internal/application/outbound"
	"user-domain/internal/domain/outport"
	"user-domain/internal/entity"
)

type orchestrator struct {
	producer      outbound.Producer
	subscriber    outbound.Subscriber
	logger        outbound.Logger
	workflowRuner outbound.WorkflowRuner
}

func (o *orchestrator) ExecuteUserUpdation(ctx context.Context, revertUser, newUser *entity.User) error {
	o.logger.Info("Starting user update orchestration for user ID: %s", newUser.ID)
	updateWorkflow := workflow.NewUpdationUserWorkflow(o.producer, o.subscriber, revertUser, newUser)
	err := o.workflowRuner.Execute(ctx, updateWorkflow.Run)
	if err != nil {
		o.logger.Error("User update workflow failed: %v", err)
		o.logExecutionTrace(updateWorkflow)
		return err
	}
	o.logger.Info("User update orchestration completed successfully for user ID: %s", newUser.ID)
	o.logExecutionTrace(updateWorkflow)
	return nil
}

func (o *orchestrator) logExecutionTrace(wf workflow.Workflow) {
	logs := wf.GetExecutionLogs()
	o.logger.Info("Workflow final state: %s", wf.GetState())

	for _, log := range logs {
		msg := fmt.Sprintf("Step %d [%s]: %s", log.StepIndex+1, log.StepName, log.State)
		if log.Error != nil {
			msg += fmt.Sprintf(" - Error: %v", log.Error)
			o.logger.Error("%s", msg)
		} else {
			o.logger.Info("%s", msg)
		}
	}
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
