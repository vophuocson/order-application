package logger

import (
	"context"
	"fmt"
	"os"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	requestID = "request_id"
)

type logger struct {
	zap *zap.Logger
}

func (l *logger) Debug(message string, fs domainoutport.LogFields) {
	l.zap.Debug(message, convertToZapField(fs)...)
}

func (l *logger) Debugf(format string, a ...any) {
	l.zap.Debug(fmt.Sprintf(format, a...))
}

func (l *logger) Info(message string, fs domainoutport.LogFields) {
	l.zap.Info(message, convertToZapField(fs)...)
}

func (l *logger) Infof(format string, a ...any) {
	l.zap.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Warn(message string, fs domainoutport.LogFields) {
	l.zap.Warn(message, convertToZapField(fs)...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.zap.Warn(fmt.Sprintf(format, a...))
}

func (l *logger) Error(message string, fs domainoutport.LogFields) {
	l.zap.Error(message, convertToZapField(fs)...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.zap.Error(fmt.Sprintf(format, a...))
}

func (l *logger) WithContext(ctx context.Context) applicationoutbound.Logger {
	rID, ok := ctx.Value(middleware.GetReqID(ctx)).(string)
	if !ok || rID == "" {
		return l
	}
	return &logger{zap: l.zap.With(zap.String(requestID, rID))}
}

func convertToZapField(fs domainoutport.LogFields) []zap.Field {
	var result []zap.Field
	for key, f := range fs {
		result = append(result, zap.Any(key, f))
	}
	return result
}

func NewLogger() applicationoutbound.Logger {
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
		zap: zap.New(core),
	}
	return &logger
}
