package applicationlogger

import (
	"context"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"
)

type logger struct {
	outBound applicationoutbound.Logger
}

func (l *logger) Debug(message string, fs domainoutport.LogFields) {
	l.outBound.Debug(message, fs)
}

func (l *logger) Debugf(format string, a ...any) {
	l.outBound.Debugf(format, a)
}

func (l *logger) Info(message string, fs domainoutport.LogFields) {
	l.outBound.Info(message, fs)
}

func (l *logger) Infof(format string, a ...any) {
	l.outBound.Infof(format, a)
}

func (l *logger) Warn(message string, fs domainoutport.LogFields) {
	l.outBound.Warn(message, fs)
}

func (l *logger) Warnf(format string, a ...any) {
	l.outBound.Warnf(format, a)
}

func (l *logger) Error(message string, fs domainoutport.LogFields) {
	l.outBound.Error(message, fs)
}

func (l *logger) Errorf(format string, a ...any) {
	l.outBound.Errorf(format, a)
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
