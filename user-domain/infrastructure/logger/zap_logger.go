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

const requestID = "request_id"

type logger struct {
	zap *zap.Logger
}

// ---- Core Logging ----

func (l *logger) log(level zapcore.Level, msg string, args ...any) {
	var fields []zap.Field
	var formatted string

	// Detect if last argument is LogFields
	if len(args) > 0 {
		if fs, ok := args[len(args)-1].(domainoutport.LogFields); ok {
			fields = convertToZapField(fs)
			args = args[:len(args)-1]
		}
	}

	// Format message if args still exist
	if len(args) > 0 {
		formatted = fmt.Sprintf(msg, args...)
	} else {
		formatted = msg
	}

	switch level {
	case zapcore.DebugLevel:
		l.zap.Debug(formatted, fields...)
	case zapcore.InfoLevel:
		l.zap.Info(formatted, fields...)
	case zapcore.WarnLevel:
		l.zap.Warn(formatted, fields...)
	case zapcore.ErrorLevel:
		l.zap.Error(formatted, fields...)
	}
}

func (l *logger) Debug(msg string, args ...any) { l.log(zapcore.DebugLevel, msg, args...) }
func (l *logger) Info(msg string, args ...any)  { l.log(zapcore.InfoLevel, msg, args...) }
func (l *logger) Warn(msg string, args ...any)  { l.log(zapcore.WarnLevel, msg, args...) }
func (l *logger) Error(msg string, args ...any) { l.log(zapcore.ErrorLevel, msg, args...) }

// ---- Context binding ----

func (l *logger) WithContext(ctx context.Context) applicationoutbound.Logger {
	rID := middleware.GetReqID(ctx)
	if rID == "" {
		return l
	}
	return &logger{zap: l.zap.With(zap.String(requestID, rID))}
}

// ---- Helpers ----

func convertToZapField(fs domainoutport.LogFields) []zap.Field {
	fields := make([]zap.Field, 0, len(fs))
	for k, v := range fs {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// ---- Lifecycle ----

func (l *logger) Sync() {
	_ = l.zap.Sync()
}

// ---- Factory ----

func NewLogger() applicationoutbound.Logger {
	env := os.Getenv("ENV")

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder
	var level zapcore.Level

	if env == "prod" {
		encoderCfg = zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			MessageKey:    "msg",
			NameKey:       "logger",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		}
		encoder = zapcore.NewJSONEncoder(encoderCfg)
		level = zap.InfoLevel
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &logger{zap: z}
}
