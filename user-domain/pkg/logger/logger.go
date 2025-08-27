package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	outport "user-domain/internal/outport"
)

type logger struct {
	logger *zap.Logger
}

func (l *logger) Debug(message string, fs outport.LogFields) {
	l.logger.Debug(message, convertToZapField(fs)...)
}

func (l *logger) Debugf(format string, a ...any) {
	l.logger.Debug(fmt.Sprintf(format, a...))
}

func (l *logger) Info(message string, fs outport.LogFields) {
	l.logger.Info(message, convertToZapField(fs)...)
}

func (l *logger) Infof(format string, a ...any) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Warn(message string, fs outport.LogFields) {
	l.logger.Warn(message, convertToZapField(fs)...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.logger.Warn(fmt.Sprintf(format, a...))
}

func (l *logger) Error(message string, fs outport.LogFields) {
	l.logger.Error(message, convertToZapField(fs)...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.logger.Error(fmt.Sprintf(format, a...))
}

func convertToZapField(fs outport.LogFields) []zap.Field {
	var result []zap.Field
	for key, f := range fs {
		result = append(result, zap.Any(key, f))
	}
	return result
}

func NewLogger() outport.Logger {
	configLog := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.EpochNanosTimeEncoder,
	}
	jsonEndcoder := zapcore.NewJSONEncoder(configLog)
	core := zapcore.NewCore(jsonEndcoder, os.Stdout, zap.DebugLevel)
	logger := logger{
		logger: zap.New(core, nil),
	}
	zap.NewExample()
	return &logger
}
