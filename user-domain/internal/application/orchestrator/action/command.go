package activity

import "context"

type Command interface {
	ExecutePending(ctx context.Context) error
	Verify(ctx context.Context) error
	Approve(ctx context.Context) error
	Compensate(ctx context.Context) error
	Name() string
}

type VerificationResponse struct {
	ServiceName string
	Accepted    bool
	Message     string
	Error       error
}
