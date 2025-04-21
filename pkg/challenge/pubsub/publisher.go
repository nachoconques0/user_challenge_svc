package pubsub

import (
	"context"

	"github.com/google/uuid"
)

type Publisher interface {
	Emit(ctx context.Context, eventID uuid.UUID, eventType string, payload interface{}) error
}
