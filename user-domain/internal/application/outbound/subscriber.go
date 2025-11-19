package outbound

import "context"

type Subscriber interface {
	Consume(ctx context.Context, channel string) ([]byte, error)
}
