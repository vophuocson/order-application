package logger

import (
	"context"
	"user-domain/internal/application/outbound"
	"user-domain/internal/domain/outport"
)

type logger struct {
	outBound outbound.Logger
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

func (l *logger) WithContext(ctx context.Context) outport.Logger {
	l.outBound.WithContext(ctx)
	return l
}

func NewLogger(o outbound.Logger) outport.Logger {
	return &logger{
		outBound: o,
	}
}
