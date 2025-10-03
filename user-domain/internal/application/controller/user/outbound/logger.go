package userapioutbound

import "context"

type FieldMap interface {
	GetKey() string
}

type Logger interface {
	Debug(message string, f ...FieldMap)
	Debugf(format string, a ...any)
	Info(message string, f ...FieldMap)
	Infof(format string, a ...any)
	Warn(message string, f ...FieldMap)
	Warnf(format string, a ...any)
	Error(message string, f ...FieldMap)
	Errorf(format string, a ...any)
	Sync() error
	WithContext(ctx context.Context) Logger
}
