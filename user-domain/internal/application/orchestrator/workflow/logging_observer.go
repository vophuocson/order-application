package workflow

import (
	"fmt"
	"user-domain/internal/application/outbound"
)

// LoggingObserver implements WorkflowObserver for logging workflow events
type LoggingObserver struct {
	logger outbound.Logger
}

// NewLoggingObserver creates a new logging observer
func NewLoggingObserver(logger outbound.Logger) WorkflowObserver {
	return &LoggingObserver{
		logger: logger,
	}
}

// OnEvent handles workflow events by logging them
func (o *LoggingObserver) OnEvent(event *WorkflowEvent) {
	message := o.formatEventMessage(event)

	if event.Error != nil {
		o.logger.Error("%s", message)
	} else {
		o.logger.Info("%s", message)
	}
}

// formatEventMessage formats an event into a log message
func (o *LoggingObserver) formatEventMessage(event *WorkflowEvent) string {
	timestamp := event.Timestamp.Format("2006-01-02 15:04:05.000")

	var message string
	if event.Phase != "" {
		message = fmt.Sprintf("[%s] Phase: %s | Step: %s (idx:%d) | Type: %s",
			timestamp, event.Phase, event.StepName, event.StepIndex, event.Type)
	} else {
		message = fmt.Sprintf("[%s] Step: %s (idx:%d) | Type: %s",
			timestamp, event.StepName, event.StepIndex, event.Type)
	}

	if event.State != "" {
		message += fmt.Sprintf(" | State: %s", event.State)
	}

	if event.Error != nil {
		message += fmt.Sprintf(" | Error: %v", event.Error)
	}

	return message
}
