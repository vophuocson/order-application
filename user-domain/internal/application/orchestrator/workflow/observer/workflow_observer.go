package observer

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

// WorkflowContext contains tracing information for distributed systems
type WorkflowContext struct {
	WorkflowID   string // Unique workflow execution ID
	WorkflowType string // Type of workflow (e.g., "user_updation")
	EntityID     string // ID of entity being processed (e.g., user ID)
}

// WorkflowEvent represents an event that occurs during workflow execution
type WorkflowEvent struct {
	Type      EventType
	StepName  string
	StepIndex int
	Phase     string
	Error     error
	Timestamp time.Time
	Duration  time.Duration // Duration of the operation
	Context   *WorkflowContext
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

// WithContext adds workflow context to the event (for tracing)
func (e *WorkflowEvent) WithContext(ctx *WorkflowContext) *WorkflowEvent {
	e.Context = ctx
	return e
}

// WithDuration sets the duration for the event
func (e *WorkflowEvent) WithDuration(duration time.Duration) *WorkflowEvent {
	e.Duration = duration
	return e
}

// AddMetadata adds custom metadata to the event
func (e *WorkflowEvent) AddMetadata(key string, value interface{}) *WorkflowEvent {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WorkflowObserver defines the interface for workflow event observers
type WorkflowObserver interface {
	// OnEvent is called when a workflow event occurs
	OnEvent(event *WorkflowEvent)
}
