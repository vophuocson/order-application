package applicationoutbound

import (
	"context"
)

type Logger interface {
	Debug(format string, a ...any)
	Info(format string, a ...any)
	Warn(format string, a ...any)
	Error(format string, a ...any)
	WithContext(ctx context.Context) Logger
}
