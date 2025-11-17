package action

import "context"

type Execute func(ctx context.Context) error
type Verify func(ctx context.Context) error
type Approve func(ctx context.Context) error
type Ran func() bool

type Execution interface {
	Execute(ctx context.Context) error
	Ran() bool
}
type Verification interface {
	Ran() bool
	Verify(ctx context.Context) error
}
type Compensation interface {
	Compensate(ctx context.Context) error
	Ran() bool
}
type Approval interface {
	Approve(ctx context.Context) error
	Ran() bool
}

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
