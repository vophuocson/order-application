package inport

import "context"

type Consumer interface {
	Consume(ctx context.Context, channelName string) (any, error)
}
