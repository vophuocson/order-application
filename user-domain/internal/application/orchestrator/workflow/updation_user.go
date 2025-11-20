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

// UserUpdationActivity handles the execution of saga activities
// It contains execute, verify, compensate, and approve methods
type UserUpdationActivity struct {
	executedSteps []int
	executionLogs []*SagaExecuteLog
	executions    []command.Execution
	verifications []command.Verification
	compensations []command.Compensation
	approvals     []command.Approval
}

// NewUserUpdationActivity creates a new activity instance
func NewUserUpdationActivity() *UserUpdationActivity {
	return &UserUpdationActivity{
		executedSteps: make([]int, 0),
		executionLogs: make([]*SagaExecuteLog, 0),
		executions:    make([]command.Execution, 0),
		verifications: make([]command.Verification, 0),
		compensations: make([]command.Compensation, 0),
		approvals:     make([]command.Approval, 0),
	}
}

// addExecution adds an execution command to the activity
func (a *UserUpdationActivity) addExecution(exec command.Execution) {
	a.executions = append(a.executions, exec)
}

// addVerification adds a verification command to the activity
func (a *UserUpdationActivity) addVerification(verify command.Verification) {
	a.verifications = append(a.verifications, verify)
}

// addCompensation adds a compensation command to the activity
func (a *UserUpdationActivity) addCompensation(compensate command.Compensation) {
	a.compensations = append(a.compensations, compensate)
}

// addApproval adds an approval command to the activity
func (a *UserUpdationActivity) addApproval(approve command.Approval) {
	a.approvals = append(a.approvals, approve)
}

func (a *UserUpdationActivity) logStep(stepName string, stepIndex int, state string, err error) {
	a.executionLogs = append(a.executionLogs, &SagaExecuteLog{
		StepName:  stepName,
		StepIndex: stepIndex,
		State:     state,
		Error:     err,
	})
}

// execute runs all execution commands in parallel
func (a *UserUpdationActivity) Execute(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.executions))
	for idx, e := range a.executions {
		if e.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, e command.Execution) {
			defer waitGroup.Done()
			a.logStep(e.Name(), idx, "command", nil)
			err := e.Execute(ctx)
			if err != nil {
				a.logStep(e.Name(), idx, "command: Failed", err)
				errChan <- fmt.Errorf("command %s failed: %w", e.Name(), err)
				return
			}
			a.logStep(e.Name(), idx, "execute: Sent", nil)
			a.executedSteps = append(a.executedSteps, idx)
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
func (a *UserUpdationActivity) Verify(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	var errChan = make(chan error, len(a.verifications))
	for idx, v := range a.verifications {
		if v.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, v command.Verification) {
			defer waitGroup.Done()
			a.logStep(v.Name(), idx, "verify", nil)
			err := v.Verify(ctx)
			if err != nil {
				a.logStep(v.Name(), idx, "verify: error", err)
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
func (a *UserUpdationActivity) Compensate(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.compensations))
	for idx, c := range a.compensations {
		if c.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, c command.Compensation) {
			defer waitGroup.Done()
			a.logStep(c.Name(), idx, "compensate", nil)
			err := c.Compensate(ctx)
			if err != nil {
				a.logStep(c.Name(), idx, "compensate: error", err)
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
func (a *UserUpdationActivity) Approve(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.approvals))
	for idx, approval := range a.approvals {
		if approval.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, approval command.Approval) {
			defer waitGroup.Done()
			a.logStep(approval.Name(), idx, "approve", nil)
			err := approval.Approve(ctx)
			if err != nil {
				a.logStep(approval.Name(), idx, "approve: error", err)
				errChan <- fmt.Errorf("approve %s failed: %w", approval.Name(), err)
				return
			}
		}(idx, approval)
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

// UserUpdationWorkflow implements the Workflow interface for user update operations
// It follows the Saga pattern with Execute, Verify, Approve, and Compensate phases
// This is the parent object that orchestrates the workflow using UserUpdationActivity
type UserUpdationWorkflow struct {
	state     SagaState
	lastError error
	activity  *UserUpdationActivity
}

func (w *UserUpdationWorkflow) logStep(stepName string, stepIndex int, state string, err error) {
	w.activity.logStep(stepName, stepIndex, state, err)
}

func (w *UserUpdationWorkflow) GetState() string {
	return w.state.String()
}

func (w *UserUpdationWorkflow) GetExecutionLogs() []*SagaExecuteLog {
	return w.activity.executionLogs
}

func (w *UserUpdationWorkflow) GetLastError() error {
	return w.lastError
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
	err := w.activity.Execute(ctx)
	if err != nil {
		w.state = SagaStateFailed
		w.lastError = err
		w.logStep("All Commands", 0, "Phase 1: Execute Failed", err)
		compensateErr := w.activity.Compensate(ctx)
		if compensateErr != nil {
			return fmt.Errorf("command failed and compensation also failed: execution error: %v, compensation error: %v", err, compensateErr)
		}
		return fmt.Errorf("execution phase failed: %w", err)
	}

	// Phase 2: Verify - Check if all services accepted pending data (PARALLEL)
	w.logStep("All Commands", 0, "Phase 2: Verifying", nil)
	if err := w.activity.Verify(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Compensation: Rollback all pending data
		w.logStep("All Commands", 0, "Compensating (Verification Failed)", err)
		if compErr := w.activity.Compensate(ctx); compErr != nil {
			return fmt.Errorf("verification failed and compensation also failed: verification error: %v, compensation error: %v", err, compErr)
		}
		return fmt.Errorf("verification phase failed: %w", err)
	}

	// Phase 3: Approve - Commit the changes across all services (PARALLEL)
	w.logStep("All Commands", 0, "Phase 3: Approving", nil)
	if err := w.activity.Approve(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		w.logStep("All Commands", 0, "Compensating (Approve Failed)", err)
		if compErr := w.activity.Compensate(ctx); compErr != nil {
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
	// Create the activity (child object)
	activity := NewUserUpdationActivity()

	// User update approval
	userApproval := action.NewUserUpdateApproval(producer, newUser.ID)
	activity.addApproval(userApproval)

	// User update compensation (rollback)
	userCompensate := action.NewUserUpdateCompensation(producer, oldUser)
	activity.addCompensation(userCompensate)

	// Payment update execution
	paymentExecution := action.NewPaymentUpdateExecution(producer, newUser)
	activity.addExecution(paymentExecution)

	// Payment update compensation (rollback)
	paymentCompensation := action.NewPaymentUpdateCompensation(producer, oldUser)
	activity.addCompensation(paymentCompensation)

	// Payment update verification
	paymentVerification := action.NewPaymentUpdateVerification(subscriber)
	activity.addVerification(paymentVerification)

	// Payment update approval
	paymentApproval := action.NewPaymentUpdateApproval(producer, newUser.ID)
	activity.addApproval(paymentApproval)

	// Create the workflow (parent object)
	workflow := &UserUpdationWorkflow{
		state:    SagaStateInitial,
		activity: activity,
	}

	return workflow
}
