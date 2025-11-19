package outbound

import "context"

type WorkflowRuner interface {
	Execute(ctx context.Context, workflowID string, taskQueueName string, workflow func(ctx context.Context) error) error
}
