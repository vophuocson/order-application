package workflow

import (
	"context"
	"fmt"
	"sync"
	"user-domain/internal/application/orchestrator/action"
	"user-domain/internal/application/orchestrator/command"
	"user-domain/internal/application/outbound"
	"user-domain/internal/entity"
)

// SagaState represents the current state of a saga workflow
type SagaState int

const (
	SagaStateInitial SagaState = iota
	SagaStateRunning
	SagaStateCompleted
	SagaStateFailed
	SagaStateCompensating
	SagaStateCompensated
)

func (s SagaState) String() string {
	switch s {
	case SagaStateInitial:
		return "Initial"
	case SagaStateRunning:
		return "Running"
	case SagaStateCompleted:
		return "Completed"
	case SagaStateFailed:
		return "Failed"
	case SagaStateCompensating:
		return "Compensating"
	case SagaStateCompensated:
		return "Compensated"
	default:
		return "Unknown"
	}
}

// SagaExecuteLog represents a log entry for a saga step execution
type SagaExecuteLog struct {
	StepName  string
	StepIndex int
	State     string
	Error     error
}

// UserUpdationWorkflow implements the Workflow interface for user update operations
// It follows the Saga pattern with Execute, Verify, Approve, and Compensate phases
type UserUpdationWorkflow struct {
	state         SagaState
	executedSteps []int
	executionLogs []*SagaExecuteLog
	lastError     error
	executions    []command.Execution
	verifications []command.Verification
	compensations []command.Compensation
	approvals     []command.Approval
}

func (w *UserUpdationWorkflow) logStep(stepName string, stepIndex int, state string, err error) {
	w.executionLogs = append(w.executionLogs, &SagaExecuteLog{
		StepName:  stepName,
		StepIndex: stepIndex,
		State:     state,
		Error:     err,
	})
}

func (w *UserUpdationWorkflow) GetState() string {
	return w.state.String()
}

func (w *UserUpdationWorkflow) GetExecutionLogs() []*SagaExecuteLog {
	return w.executionLogs
}

func (w *UserUpdationWorkflow) GetLastError() error {
	return w.lastError
}

// execute runs all execution commands in parallel
func (w *UserUpdationWorkflow) execute(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.executions))
	for idx, e := range w.executions {
		if e.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, e command.Execution) {
			defer waitGroup.Done()
			w.logStep(e.Name(), idx, "command", nil)
			err := e.Execute(ctx)
			if err != nil {
				w.logStep(e.Name(), idx, "command: Failed", err)
				errChan <- fmt.Errorf("command %s failed: %w", e.Name(), err)
				return
			}
			w.logStep(e.Name(), idx, "execute: Sent", nil)
			w.executedSteps = append(w.executedSteps, idx)
		}(idx, e)
	}
	waitGroup.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// verify runs all verification commands in parallel
func (w *UserUpdationWorkflow) verify(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	var errChan = make(chan error, len(w.verifications))
	for idx, v := range w.verifications {
		if v.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, v command.Verification) {
			defer waitGroup.Done()
			w.logStep(v.Name(), idx, "verify", nil)
			err := v.Verify(ctx)
			if err != nil {
				w.logStep(v.Name(), idx, "verify: error", err)
				errChan <- fmt.Errorf("verify %s failed: %w", v.Name(), err)
				return
			}
		}(idx, v)
	}
	waitGroup.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// compensate runs all compensation commands in parallel (for rollback)
func (w *UserUpdationWorkflow) compensate(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.compensations))
	for idx, c := range w.compensations {
		if c.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, c command.Compensation) {
			defer waitGroup.Done()
			w.logStep(c.Name(), idx, "compensate", nil)
			err := c.Compensate(ctx)
			if err != nil {
				w.logStep(c.Name(), idx, "compensate: error", err)
				errChan <- fmt.Errorf("compensate %s failed: %w", c.Name(), err)
				return
			}
		}(idx, c)
	}
	waitGroup.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// approve runs all approval commands in parallel
func (w *UserUpdationWorkflow) approve(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.approvals))
	for idx, a := range w.approvals {
		if a.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, a command.Approval) {
			defer waitGroup.Done()
			w.logStep(a.Name(), idx, "approve", nil)
			err := a.Approve(ctx)
			if err != nil {
				w.logStep(a.Name(), idx, "approve: error", err)
				errChan <- fmt.Errorf("approve %s failed: %w", a.Name(), err)
				return
			}
		}(idx, a)
	}
	waitGroup.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// Run executes the complete saga workflow with all phases
// Phase 1: Execute - Send pending requests to all services
// Phase 2: Verify - Wait for all services to accept pending data
// Phase 3: Approve - Commit the changes across all services
// If any phase fails, compensate (rollback) all changes
func (w *UserUpdationWorkflow) Run(ctx context.Context) error {
	// Phase 1: Execute - Send pending requests to all services (PARALLEL)
	w.logStep("All Commands", 0, "Phase 1: Executing Pending", nil)
	w.state = SagaStateRunning
	err := w.execute(ctx)
	if err != nil {
		w.state = SagaStateFailed
		w.lastError = err
		w.logStep("All Commands", 0, "Phase 1: Execute Failed", err)
		compensateErr := w.compensate(ctx)
		if compensateErr != nil {
			return fmt.Errorf("command failed and compensation also failed: execution error: %v, compensation error: %v", err, compensateErr)
		}
		return fmt.Errorf("execution phase failed: %w", err)
	}

	// Phase 2: Verify - Check if all services accepted pending data (PARALLEL)
	w.logStep("All Commands", 0, "Phase 2: Verifying", nil)
	if err := w.verify(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Compensation: Rollback all pending data
		w.logStep("All Commands", 0, "Compensating (Verification Failed)", err)
		if compErr := w.compensate(ctx); compErr != nil {
			return fmt.Errorf("verification failed and compensation also failed: verification error: %v, compensation error: %v", err, compErr)
		}
		return fmt.Errorf("verification phase failed: %w", err)
	}

	// Phase 3: Approve - Commit the changes across all services (PARALLEL)
	w.logStep("All Commands", 0, "Phase 3: Approving", nil)
	if err := w.approve(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		w.logStep("All Commands", 0, "Compensating (Approve Failed)", err)
		if compErr := w.compensate(ctx); compErr != nil {
			return fmt.Errorf("approve failed and compensation also failed: approve error: %v, compensation error: %v", err, compErr)
		}

		return fmt.Errorf("approve phase failed: %w", err)
	}

	w.state = SagaStateCompleted
	w.logStep("All Commands", 0, "Completed Successfully", nil)
	return nil
}

// NewUpdationUserWorkflow creates a new user update workflow with all necessary commands
// This workflow coordinates user updates across multiple services using the Saga pattern
func NewUpdationUserWorkflow(producer outbound.Producer, subscriber outbound.Subscriber, oldUser, newUser *entity.User) Workflow {
	workflow := &UserUpdationWorkflow{
		state: SagaStateInitial,
	}

	// User update approval
	userApproval := action.NewUserUpdateApproval(producer, newUser.ID)
	workflow.approvals = append(workflow.approvals, userApproval)

	// User update compensation (rollback)
	userCompensate := action.NewUserUpdateCompensation(producer, oldUser)
	workflow.compensations = append(workflow.compensations, userCompensate)

	// Payment update execution
	paymentExecution := action.NewPaymentUpdateExecution(producer, newUser)
	workflow.executions = append(workflow.executions, paymentExecution)

	// Payment update compensation (rollback)
	paymentCompensation := action.NewPaymentUpdateCompensation(producer, oldUser)
	workflow.compensations = append(workflow.compensations, paymentCompensation)

	// Payment update verification
	paymentVerification := action.NewPaymentUpdateVerification(subscriber)
	workflow.verifications = append(workflow.verifications, paymentVerification)

	// Payment update approval
	paymentApproval := action.NewPaymentUpdateApproval(producer, newUser.ID)
	workflow.approvals = append(workflow.approvals, paymentApproval)

	return workflow
}
