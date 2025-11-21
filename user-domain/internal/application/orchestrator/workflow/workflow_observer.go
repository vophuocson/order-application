package workflow

import "time"

// EventType represents the type of workflow event
type EventType string

const (
	EventTypeExecuteStart      EventType = "execute.start"
	EventTypeExecuteSuccess    EventType = "execute.success"
	EventTypeExecuteFailed     EventType = "execute.failed"
	EventTypeVerifyStart       EventType = "verify.start"
	EventTypeVerifySuccess     EventType = "verify.success"
	EventTypeVerifyFailed      EventType = "verify.failed"
	EventTypeApproveStart      EventType = "approve.start"
	EventTypeApproveSuccess    EventType = "approve.success"
	EventTypeApproveFailed     EventType = "approve.failed"
	EventTypeCompensateStart   EventType = "compensate.start"
	EventTypeCompensateSuccess EventType = "compensate.success"
	EventTypeCompensateFailed  EventType = "compensate.failed"
	EventTypePhaseStart        EventType = "phase.start"
	EventTypePhaseComplete     EventType = "phase.complete"
	EventTypeWorkflowComplete  EventType = "workflow.complete"
)

// WorkflowEvent represents an event that occurs during workflow execution
type WorkflowEvent struct {
	Type      EventType
	StepName  string
	StepIndex int
	Phase     string
	State     string
	Error     error
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// NewWorkflowEvent creates a new workflow event
func NewWorkflowEvent(eventType EventType, stepName string, stepIndex int, phase string, err error) *WorkflowEvent {
	return &WorkflowEvent{
		Type:      eventType,
		StepName:  stepName,
		StepIndex: stepIndex,
		Phase:     phase,
		Error:     err,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// WorkflowObserver defines the interface for workflow event observers
type WorkflowObserver interface {
	// OnEvent is called when a workflow event occurs
	OnEvent(event *WorkflowEvent)
}
