package outbound

import "context"

type WorkflowRuner interface {
	Execute(ctx context.Context, workflow func(ctx context.Context) error) error
}
