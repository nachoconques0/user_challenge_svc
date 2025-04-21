package local

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"

	eventUser "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user/event"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/pubsub"
)

type Bus struct {
	db          *gorm.DB
	mu          sync.RWMutex
	subscribers map[string][]pubsub.HandlerFunc
}

func NewBus(db *gorm.DB) *Bus {
	return &Bus{
		db:          db,
		subscribers: make(map[string][]pubsub.HandlerFunc),
	}
}

func (b *Bus) Publish(ctx context.Context, eventID string, eventType string, payload any) error {
	tx := b.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Model(&eventUser.User{}).
		Where("id = ?", eventID).
		Update("published", true).Error
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("event_id", eventID).Msg("failed to update published event")
		return err
	}

	b.mu.RLock()
	handlers := b.subscribers[eventType]
	b.mu.RUnlock()

	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("failed to commit after event publish")
		return err
	}

	log.Info().
		Str("event_type", eventType).
		Str("event_id", eventID).
		Msg("event published and marked as published in DB")

	for _, handler := range handlers {
		go handler(ctx, payload)
	}

	return nil
}

func (b *Bus) Subscribe(event string, handler pubsub.HandlerFunc) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[event] = append(b.subscribers[event], handler)
	return nil
}
