package workflow

import "context"

// Workflow represents a saga workflow that can be executed
type Workflow interface {
	// Run executes the complete workflow including execute, verify, approve and compensate on failure
	Run(ctx context.Context) error
	// GetState returns the current state of the workflow
	GetState() string
	// GetExecutionLogs returns all execution logs
	GetExecutionLogs() []*SagaExecuteLog
	// GetLastError returns the last error that occurred
	GetLastError() error
}

