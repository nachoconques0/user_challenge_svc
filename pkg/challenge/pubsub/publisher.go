package pubsub

import "context"

type Publisher interface {
	Publish(ctx context.Context, eventID string, eventType string, payload interface{}) error
}
