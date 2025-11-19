package outbound

import "context"

type Producer interface {
	Push(ctx context.Context, channel string, data []byte) error
}
