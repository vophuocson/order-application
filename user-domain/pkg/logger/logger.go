package logger

import (
	"fmt"
	"user-domain/internal/outport"
	lg "user-domain/internal/outport"

	"go.uber.org/zap"
)

type logger struct {
	logger *zap.Logger
}

func (l *logger) Debug(message string, fs lg.LogFields) {
	l.logger.Debug(message, convertToZapField(fs)...)
}

func (l *logger) Debugf(format string, a ...any) {
	l.logger.Debug(fmt.Sprintf(format, a...))
}

func (l *logger) Info(message string, fs lg.LogFields) {
	l.logger.Info(message, convertToZapField(fs)...)
}

func (l *logger) Infof(format string, a ...any) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Warn(message string, fs lg.LogFields) {
	l.logger.Warn(message, convertToZapField(fs)...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.logger.Warn(fmt.Sprintf(format, a...))
}

func (l *logger) Error(message string, fs lg.LogFields) {
	l.logger.Error(message, convertToZapField(fs)...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.logger.Error(fmt.Sprintf(format, a...))
}

func convertToZapField(fs lg.LogFields) []zap.Field {
	var result []zap.Field
	for key, f := range fs {
		result = append(result, zap.Any(key, f))
	}
	return result
}

func NewLogger() outport.Logger {
	return &logger{}
}
