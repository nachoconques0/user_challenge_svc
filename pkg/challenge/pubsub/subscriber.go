package pubsub

import "context"

type HandlerFunc func(ctx context.Context, payload interface{})

type Subscriber interface {
	Subscribe(event string, handler HandlerFunc) error
}
