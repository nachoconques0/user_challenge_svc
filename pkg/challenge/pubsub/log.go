package pubsub

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	eventUser "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user/event"
)

type LoggerPublisher struct {
	db *gorm.DB
}

func NewLoggerPublisher(db *gorm.DB) *LoggerPublisher {
	return &LoggerPublisher{db: db}
}

// Emit logs the event and marks it as published in the DB
func (p *LoggerPublisher) Emit(_ context.Context, eventID uuid.UUID, eventType string, payload interface{}) error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Model(&eventUser.User{}).
		Where("id = ?", eventID).
		Update("published", true).Error
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("event_id", eventID.String()).Msg("failed to update published event")
		return err
	}

	// Here because it is a challeng and we are not emiting anywhere we are just loggin.
	// but in the scenario of emitin an event to a topic and it fails. We should not update in the DB
	log.Info().
		Str("event_type", eventType).
		Str("event_id", eventID.String()).
		Interface("payload", payload).
		Msg("Event emitted successfully")

	// here we roll back if there was an error.

	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Str("event_id", eventID.String()).Msg("failed to update event")
		return err
	}

	return nil
}
