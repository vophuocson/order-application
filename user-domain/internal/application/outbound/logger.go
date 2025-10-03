package applicationoutbound

import domainoutport "user-domain/internal/domain/outport"

type Logger interface {
	Debug(message string, f domainoutport.LogFields)
	Debugf(format string, a ...any)
	Info(message string, f domainoutport.LogFields)
	Infof(format string, a ...any)
	Warn(message string, f domainoutport.LogFields)
	Warnf(format string, a ...any)
	Error(message string, f domainoutport.LogFields)
	Errorf(format string, a ...any)
	// Sync() error
	// WithContext(ctx context.Context) Logger
}
