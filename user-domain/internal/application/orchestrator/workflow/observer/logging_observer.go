package observer

import (
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

	if event.Error != nil {
		o.logger.Error("workflow_log", event)
	} else {
		o.logger.Info("workflow_log", event)
	}
}
