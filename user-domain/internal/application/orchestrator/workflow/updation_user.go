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

// ActivityStep represents a single step in the saga with all its phases
// ActivityStep apply builder patter
type ActivityStep struct {
	name         string
	execution    command.Execution
	compensation command.Compensation
	verification command.Verification
	approval     command.Approval
	executed     bool // tracks if execution succeeded
}

// NewActivityStep creates a new activity step
func NewActivityStep(name string) *ActivityStep {
	return &ActivityStep{
		name:     name,
		executed: false,
	}
}

// SetExecution sets the execution command for this step
func (s *ActivityStep) SetExecution(exec command.Execution) *ActivityStep {
	s.execution = exec
	return s
}

// SetCompensation sets the compensation command for this step
func (s *ActivityStep) SetCompensation(comp command.Compensation) *ActivityStep {
	s.compensation = comp
	return s
}

// SetVerification sets the verification command for this step
func (s *ActivityStep) SetVerification(verify command.Verification) *ActivityStep {
	s.verification = verify
	return s
}

// SetApproval sets the approval command for this step
func (s *ActivityStep) SetApproval(approve command.Approval) *ActivityStep {
	s.approval = approve
	return s
}

// UserUpdationActivity handles the execution of saga activities
// It contains execute, verify, compensate, and approve methods
type UserUpdationActivity struct {
	executionLogs []*SagaExecuteLog
	steps         []*ActivityStep
	mu            sync.Mutex // protects executionLogs and executed flags
}

// NewUserUpdationActivity creates a new activity instance
func NewUserUpdationActivity() *UserUpdationActivity {
	return &UserUpdationActivity{
		executionLogs: make([]*SagaExecuteLog, 0),
		steps:         make([]*ActivityStep, 0),
	}
}

// AddStep adds an activity step to the workflow
func (a *UserUpdationActivity) AddStep(step *ActivityStep) {
	a.steps = append(a.steps, step)
}

func (a *UserUpdationActivity) logStep(stepName string, stepIndex int, state string, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.executionLogs = append(a.executionLogs, &SagaExecuteLog{
		StepName:  stepName,
		StepIndex: stepIndex,
		State:     state,
		Error:     err,
	})
}

// Execute runs all execution commands in parallel
// Only marks step as executed if execution succeeds
func (a *UserUpdationActivity) Execute(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.steps))

	for idx, step := range a.steps {
		if step.execution == nil {
			continue
		}
		if step.execution.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, step *ActivityStep) {
			defer waitGroup.Done()
			a.logStep(step.execution.Name(), idx, "execute: start", nil)
			err := step.execution.Execute(ctx)
			if err != nil {
				a.logStep(step.execution.Name(), idx, "execute: Failed", err)
				errChan <- fmt.Errorf("execution %s failed: %w", step.execution.Name(), err)
				return
			}
			a.logStep(step.execution.Name(), idx, "execute: Success", nil)
			// Mark step as executed so compensation can run if needed
			a.mu.Lock()
			step.executed = true
			a.mu.Unlock()
		}(idx, step)
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

// Verify runs all verification commands in parallel
// Only verifies steps that were successfully executed
func (a *UserUpdationActivity) Verify(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	var errChan = make(chan error, len(a.steps))

	for idx, step := range a.steps {
		// Only verify if step was executed and has verification
		if !step.executed || step.verification == nil {
			continue
		}
		if step.verification.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, step *ActivityStep) {
			defer waitGroup.Done()
			a.logStep(step.verification.Name(), idx, "verify: start", nil)
			err := step.verification.Verify(ctx)
			if err != nil {
				a.logStep(step.verification.Name(), idx, "verify: Failed", err)
				errChan <- fmt.Errorf("verification %s failed: %w", step.verification.Name(), err)
				return
			}
			a.logStep(step.verification.Name(), idx, "verify: Success", nil)
		}(idx, step)
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

// Compensate runs all compensation commands in parallel (for rollback)
// Only compensates steps that were successfully executed
func (a *UserUpdationActivity) Compensate(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.steps))

	for idx, step := range a.steps {
		// Only compensate if step was executed and has compensation
		if !step.executed || step.compensation == nil {
			continue
		}
		// check for retry
		if step.compensation.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, step *ActivityStep) {
			defer waitGroup.Done()
			a.logStep(step.compensation.Name(), idx, "compensate: start", nil)
			err := step.compensation.Compensate(ctx)
			if err != nil {
				a.logStep(step.compensation.Name(), idx, "compensate: Failed", err)
				errChan <- fmt.Errorf("compensation %s failed: %w", step.compensation.Name(), err)
				return
			}
			a.logStep(step.compensation.Name(), idx, "compensate: Success", nil)
		}(idx, step)
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

// Approve runs all approval commands in parallel
// Only approves steps that were successfully executed
func (a *UserUpdationActivity) Approve(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(a.steps))

	for idx, step := range a.steps {
		// Only approve if step was executed and has approval
		if !step.executed || step.approval == nil {
			continue
		}
		if step.approval.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, step *ActivityStep) {
			defer waitGroup.Done()
			a.logStep(step.approval.Name(), idx, "approve: start", nil)
			err := step.approval.Approve(ctx)
			if err != nil {
				a.logStep(step.approval.Name(), idx, "approve: Failed", err)
				errChan <- fmt.Errorf("approval %s failed: %w", step.approval.Name(), err)
				return
			}
			a.logStep(step.approval.Name(), idx, "approve: Success", nil)
		}(idx, step)
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

	// Step 1: User update step (no execution, only approval and compensation)
	userStep := NewActivityStep("UserUpdate").
		SetApproval(action.NewUserUpdateApproval(producer, newUser.ID)).
		SetCompensation(action.NewUserUpdateCompensation(producer, oldUser))
	// Mark as executed since there's no execution command
	userStep.executed = true
	activity.AddStep(userStep)

	// Step 2: Payment update step (full saga: execute -> verify -> approve with compensation)
	paymentStep := NewActivityStep("PaymentUpdate").
		SetExecution(action.NewPaymentUpdateExecution(producer, newUser)).
		SetCompensation(action.NewPaymentUpdateCompensation(producer, oldUser)).
		SetVerification(action.NewPaymentUpdateVerification(subscriber)).
		SetApproval(action.NewPaymentUpdateApproval(producer, newUser.ID))
	activity.AddStep(paymentStep)

	// Create the workflow (parent object)
	workflow := &UserUpdationWorkflow{
		state:    SagaStateInitial,
		activity: activity,
	}

	return workflow
}
