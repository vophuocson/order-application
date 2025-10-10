package applicationlogger

import (
	"context"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"
)

type logger struct {
	outBound applicationoutbound.Logger
}

func (l *logger) Debug(message string, fs ...any) {
	l.outBound.Debug(message, fs)
}

func (l *logger) Info(message string, fs ...any) {
	l.outBound.Info(message, fs)
}

func (l *logger) Warn(message string, fs ...any) {
	l.outBound.Warn(message, fs)
}

func (l *logger) Error(message string, fs ...any) {
	l.outBound.Error(message, fs)
}

func (l *logger) WithContext(ctx context.Context) domainoutport.Logger {
	l.outBound.WithContext(ctx)
	return l
}

func NewLogger(o applicationoutbound.Logger) domainoutport.Logger {
	return &logger{
		outBound: o,
	}
}
