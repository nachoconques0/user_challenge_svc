package pubsub

import (
	"context"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user/event"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/pubsub"
	"github.com/rs/zerolog/log"
)

func RegisterUserSubscribers(sub pubsub.Subscriber) error {
	if err := sub.Subscribe(event.UserCreated, onUserCreated); err != nil {
		return err
	}
	if err := sub.Subscribe(event.UserUpdated, onUserUpdated); err != nil {
		return err
	}
	if err := sub.Subscribe(event.UserSoftDeleted, onUserSoftDeleted); err != nil {
		return err
	}
	return nil
}

func onUserCreated(_ context.Context, payload interface{}) {
	data, ok := payload.(event.CreatedPayload)
	if !ok {
		log.Error().Msg("invalid payload type for USER_CREATED")
		return
	}
	log.Info().Str("trace_id", data.TraceID).Str("userID", data.UserID).Msg("USER_CREATED")
}

func onUserUpdated(_ context.Context, payload interface{}) {
	data, ok := payload.(event.UpdatedPayload)
	if !ok {
		log.Error().Msg("invalid payload type for USER_UPDATED")
		return
	}
	log.Info().Str("trace_id", data.TraceID).Str("nickname", data.Nickname).Msg("USER_UPDATED")
}

func onUserSoftDeleted(_ context.Context, payload interface{}) {
	data, ok := payload.(event.DeletedPayload)
	if !ok {
		log.Error().Msg("invalid payload type for USER_SOFT_DELETED")
		return
	}
	log.Info().Str("trace_id", data.TraceID).Str("user_id", data.UserID).Msg("USER_SOFT_DELETED")
}
