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
	steps         []*ActivityStep
	mu            sync.Mutex           // protects executed flags
	eventNotifier func(*WorkflowEvent) // callback to notify observers
}

// NewUserUpdationActivity creates a new activity instance
func NewUserUpdationActivity() *UserUpdationActivity {
	return &UserUpdationActivity{
		steps:         make([]*ActivityStep, 0),
		eventNotifier: func(*WorkflowEvent) {}, // no-op by default
	}
}

// SetEventNotifier sets the event notification callback
func (a *UserUpdationActivity) SetEventNotifier(notifier func(*WorkflowEvent)) {
	a.eventNotifier = notifier
}

// AddStep adds an activity step to the workflow
func (a *UserUpdationActivity) AddStep(step *ActivityStep) {
	a.steps = append(a.steps, step)
}

// notifyEvent sends an event notification
func (a *UserUpdationActivity) notifyEvent(event *WorkflowEvent) {
	if a.eventNotifier != nil {
		a.eventNotifier(event)
	}
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

			// Notify execution start
			a.notifyEvent(NewWorkflowEvent(EventTypeExecuteStart, step.execution.Name(), idx, "execute", nil))

			err := step.execution.Execute(ctx)
			if err != nil {
				// Notify execution failure
				a.notifyEvent(NewWorkflowEvent(EventTypeExecuteFailed, step.execution.Name(), idx, "execute", err))
				errChan <- fmt.Errorf("execution %s failed: %w", step.execution.Name(), err)
				return
			}

			// Notify execution success
			a.notifyEvent(NewWorkflowEvent(EventTypeExecuteSuccess, step.execution.Name(), idx, "execute", nil))

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

			// Notify verification start
			a.notifyEvent(NewWorkflowEvent(EventTypeVerifyStart, step.verification.Name(), idx, "verify", nil))

			err := step.verification.Verify(ctx)
			if err != nil {
				// Notify verification failure
				a.notifyEvent(NewWorkflowEvent(EventTypeVerifyFailed, step.verification.Name(), idx, "verify", err))
				errChan <- fmt.Errorf("verification %s failed: %w", step.verification.Name(), err)
				return
			}

			// Notify verification success
			a.notifyEvent(NewWorkflowEvent(EventTypeVerifySuccess, step.verification.Name(), idx, "verify", nil))
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

			// Notify compensation start
			a.notifyEvent(NewWorkflowEvent(EventTypeCompensateStart, step.compensation.Name(), idx, "compensate", nil))

			err := step.compensation.Compensate(ctx)
			if err != nil {
				// Notify compensation failure
				a.notifyEvent(NewWorkflowEvent(EventTypeCompensateFailed, step.compensation.Name(), idx, "compensate", err))
				errChan <- fmt.Errorf("compensation %s failed: %w", step.compensation.Name(), err)
				return
			}

			// Notify compensation success
			a.notifyEvent(NewWorkflowEvent(EventTypeCompensateSuccess, step.compensation.Name(), idx, "compensate", nil))
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

			// Notify approval start
			a.notifyEvent(NewWorkflowEvent(EventTypeApproveStart, step.approval.Name(), idx, "approve", nil))

			err := step.approval.Approve(ctx)
			if err != nil {
				// Notify approval failure
				a.notifyEvent(NewWorkflowEvent(EventTypeApproveFailed, step.approval.Name(), idx, "approve", err))
				errChan <- fmt.Errorf("approval %s failed: %w", step.approval.Name(), err)
				return
			}

			// Notify approval success
			a.notifyEvent(NewWorkflowEvent(EventTypeApproveSuccess, step.approval.Name(), idx, "approve", nil))
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
	observers []WorkflowObserver
}

// AddObserver adds an observer to the workflow
func (w *UserUpdationWorkflow) AddObserver(observer WorkflowObserver) {
	w.observers = append(w.observers, observer)
}

// notifyObservers notifies all observers of an event
func (w *UserUpdationWorkflow) notifyObservers(event *WorkflowEvent) {
	for _, observer := range w.observers {
		observer.OnEvent(event)
	}
}

func (w *UserUpdationWorkflow) GetState() string {
	return w.state.String()
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
	phaseEvent := NewWorkflowEvent(EventTypePhaseStart, "All Commands", 0, "Phase 1: Execute", nil)
	w.notifyObservers(phaseEvent)

	w.state = SagaStateRunning
	err := w.activity.Execute(ctx)
	if err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Notify phase failure
		failEvent := NewWorkflowEvent(EventTypeExecuteFailed, "All Commands", 0, "Phase 1: Execute Failed", err)
		w.notifyObservers(failEvent)

		// Compensate
		compensateErr := w.activity.Compensate(ctx)
		if compensateErr != nil {
			return fmt.Errorf("command failed and compensation also failed: execution error: %v, compensation error: %v", err, compensateErr)
		}
		return fmt.Errorf("execution phase failed: %w", err)
	}

	// Phase 2: Verify - Check if all services accepted pending data (PARALLEL)
	verifyEvent := NewWorkflowEvent(EventTypePhaseStart, "All Commands", 0, "Phase 2: Verify", nil)
	w.notifyObservers(verifyEvent)

	if err := w.activity.Verify(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Notify verification failure
		failEvent := NewWorkflowEvent(EventTypeVerifyFailed, "All Commands", 0, "Phase 2: Verification Failed", err)
		w.notifyObservers(failEvent)

		// Compensation: Rollback all pending data
		if compErr := w.activity.Compensate(ctx); compErr != nil {
			return fmt.Errorf("verification failed and compensation also failed: verification error: %v, compensation error: %v", err, compErr)
		}
		return fmt.Errorf("verification phase failed: %w", err)
	}

	// Phase 3: Approve - Commit the changes across all services (PARALLEL)
	approveEvent := NewWorkflowEvent(EventTypePhaseStart, "All Commands", 0, "Phase 3: Approve", nil)
	w.notifyObservers(approveEvent)

	if err := w.activity.Approve(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Notify approval failure
		failEvent := NewWorkflowEvent(EventTypeApproveFailed, "All Commands", 0, "Phase 3: Approval Failed", err)
		w.notifyObservers(failEvent)

		if compErr := w.activity.Compensate(ctx); compErr != nil {
			return fmt.Errorf("approve failed and compensation also failed: approve error: %v, compensation error: %v", err, compErr)
		}

		return fmt.Errorf("approve phase failed: %w", err)
	}

	w.state = SagaStateCompleted

	// Notify workflow completion
	completeEvent := NewWorkflowEvent(EventTypeWorkflowComplete, "All Commands", 0, "Workflow Completed", nil)
	w.notifyObservers(completeEvent)

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

	// Create saga log observer for backward compatibility

	// Create the workflow (parent object)
	workflow := &UserUpdationWorkflow{
		state:    SagaStateInitial,
		activity: activity,
	}

	// Wire up activity to notify workflow observers
	activity.SetEventNotifier(workflow.notifyObservers)

	return workflow
}

// NewUpdationUserWorkflowWithObservers creates a workflow with custom observers
func NewUpdationUserWorkflowWithObservers(
	producer outbound.Producer,
	subscriber outbound.Subscriber,
	oldUser, newUser *entity.User,
	observers ...WorkflowObserver,
) Workflow {
	workflow := NewUpdationUserWorkflow(producer, subscriber, oldUser, newUser).(*UserUpdationWorkflow)

	// Add additional observers
	for _, observer := range observers {
		workflow.AddObserver(observer)
	}

	return workflow
}
