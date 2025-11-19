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

type SagaExecuteLog struct {
	StepName  string
	StepIndex int
	State     string
	Error     error
}

type UserUpdationWorkflow struct {
	state         SagaState
	executedSteps []int
	executionLogs []*SagaExecuteLog
	lastError     error
	Executions    []command.Execution
	Verifications []command.Verification
	Compensations []command.Compensation
	Approvals     []command.Approval
}

// func (u *UserUpdationWorkflow) GetExecutions() []action.Execution {
// 	return u.Executions
// }

// func (u *UserUpdationWorkflow) GetVerifications() []action.Verification {
// 	return u.Verifications
// }

// func (u *UserUpdationWorkflow) GetCompensations() []action.Compensation {
// 	return u.Compensations
// }

// func (u *UserUpdationWorkflow) GetApproval() []action.Approval {
// 	return u.Approvals
// }

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

func (w *UserUpdationWorkflow) Execute(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.Executions))
	for idx, e := range w.Executions {
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

func (w *UserUpdationWorkflow) Verify(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	var errChan = make(chan error, len(w.Verifications))
	for idx, v := range w.Verifications {
		if v.Ran() {
			continue
		}
		waitGroup.Add(1)
		go func(idx int, v command.Verification) {
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

func (w *UserUpdationWorkflow) Compensate(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.Compensations))
	for idx, c := range w.Compensations {
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

func (w *UserUpdationWorkflow) Approve(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, len(w.Approvals))
	for idx, a := range w.Approvals {
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

func (w *UserUpdationWorkflow) Run(ctx context.Context) error {
	w.logStep("All Commands", 0, "Phase 1: Executing Pending", nil)
	w.state = SagaStateRunning
	err := w.Execute(ctx)
	if err != nil {
		w.state = SagaStateFailed
		w.lastError = err
		w.logStep("All Commands", 0, "Phase 1: Executing Pending", nil)
		compensateErr := w.Compensate(ctx)
		if compensateErr != nil {
			return fmt.Errorf("command failed and compensation also failed: verification error: %v, compensation error: %v", err, compensateErr)
		}
		return fmt.Errorf("pending phase failed: %w", err)
	}

	// Phase 2: Verify - Check if all services accepted pending data (PARALLEL)
	w.logStep("All Commands", 0, "Phase 2: Verifying", nil)
	if err := w.Verify(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		// Compensation: Rollback all pending data
		w.logStep("All Commands", 0, "Phase 3: Compensating (Verification Failed)", err)
		if compErr := w.Compensate(ctx); compErr != nil {
			return fmt.Errorf("verification failed and compensation also failed: verification error: %v, compensation error: %v", err, compErr)
		}
		return fmt.Errorf("verification phase failed: %w", err)
	}

	w.logStep("All Commands", 0, "Phase 3: Approving", nil)
	if err := w.Approve(ctx); err != nil {
		w.state = SagaStateFailed
		w.lastError = err

		w.logStep("All Commands", 0, "Phase 3: Compensating (Approve Failed)", err)
		if compErr := w.Compensate(ctx); compErr != nil {
			return fmt.Errorf("approve failed and compensation also failed: approve error: %v, compensation error: %v", err, compErr)
		}

		return fmt.Errorf("approve phase failed: %w", err)
	}

	w.state = SagaStateCompleted
	w.logStep("All Commands", 0, "Completed Successfully", nil)
	return nil

}

func NewUpdationUserWorkflow(producer outbound.Producer, subscirber outbound.Subscriber, oldUser, newUser *entity.User) *UserUpdationWorkflow {
	workflow := UserUpdationWorkflow{}
	userApproval := action.NewUserUpdateApproval(producer, newUser.ID)
	workflow.Approvals = append(workflow.Approvals, userApproval)

	userCompensate := action.NewUserUpdateCompensation(producer, oldUser)
	workflow.Compensations = append(workflow.Compensations, userCompensate)

	paymentExecution := action.NewPaymentUpdateExecution(producer, newUser)
	workflow.Executions = append(workflow.Executions, paymentExecution)

	paymenCompensation := action.NewPaymentUpdateCompensation(producer, oldUser)
	workflow.Compensations = append(workflow.Compensations, paymenCompensation)

	paymenVerification := action.NewPaymentUpdateVerification(subscirber)
	workflow.Verifications = append(workflow.Verifications, paymenVerification)

	paymenApproval := action.NewPaymentUpdateApproval(producer, newUser.ID)
	workflow.Approvals = append(workflow.Approvals, paymenApproval)

	return &workflow
}
