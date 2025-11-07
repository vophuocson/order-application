package outport

import "context"

type Producer interface {
	Push(ctx context.Context, channelName string, data any) error
}
