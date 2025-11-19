package orchestrator

import (
	"context"
	"user-domain/internal/application/outbound"

	"go.temporal.io/sdk/client"
)

type orchestrator struct {
	client client.Client
}

func NewTemporalClient(client client.Client) outbound.WorkflowRuner {
	return &orchestrator{client: client}
}

func (o *orchestrator) Execute(ctx context.Context, workflowID string, taskQueueName string, workflow func(ctx context.Context) error) error {
	option := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueueName,
	}
	we, err := o.client.ExecuteWorkflow(ctx, option, workflow)
	if err != nil {
		return err
	}

	var errR error
	err = we.Get(ctx, errR)
	if err != nil {
		return err
	}

	return errR
}
