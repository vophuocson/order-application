package action

import "context"

type Execute func(ctx context.Context) error
type Verify func(ctx context.Context) error
type Approve func(ctx context.Context) error
type Ran func() bool

type Execution interface {
	Execute(ctx context.Context) error
	Ran() bool
	Name() string
}
type Verification interface {
	Ran() bool
	Verify(ctx context.Context) error
	Name() string
}
type Compensation interface {
	Compensate(ctx context.Context) error
	Ran() bool
	Name() string
}
type Approval interface {
	Approve(ctx context.Context) error
	Name() string
	Ran() bool
}

type VerificationResponse struct {
	ServiceName string
	Accepted    bool
	Message     string
	Error       error
}
