package outport

type LogFields map[string]interface{}

type Logger interface {
	Debug(message string, f LogFields)
	Debugf(format string, a ...any)
	Info(message string, f LogFields)
	Infof(format string, a ...any)
	Warn(message string, f LogFields)
	Warnf(format string, a ...any)
	Error(message string, f LogFields)
	Errorf(format string, a ...any)
	// Sync() error
	// WithContext(ctx context.Context) Logger
}
