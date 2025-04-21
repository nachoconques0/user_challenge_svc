package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dbInstance "github.com/nachoconques0/user_challenge_svc/pkg/challenge/db"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/env"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user/event"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/repo"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/pubsub"
)

type aggregate struct {
	DB        *gorm.DB
	TestTx    bool
	publisher pubsub.Publisher
}

type Aggregate interface {
	Create(ctx context.Context, u *user.Entity) (*user.Entity, error)
	Update(ctx context.Context, u *user.Entity) (*user.Entity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, country string, page, limit int) ([]user.Entity, error)
}

const (
	// ErrMissingDB used when DB is nil
	ErrMissingDB = "Aggregate is missing DB connection"
	// ErrMissingTestEnv when test env is missing
	ErrMissingTestEnv = "DB connection can only be a TX when ENV == env.Test"
)

// New returns a new User aggregate
func New(db *gorm.DB, e string, pub pubsub.Publisher) (aggregate, error) {
	a := aggregate{
		DB:        db,
		publisher: pub,
	}

	switch {
	case db == nil:
		return a, errors.New(ErrMissingDB)
	case dbInstance.IsTransaction(db) && !env.IsTest(e):
		return a, errors.New(ErrMissingTestEnv)
	case dbInstance.IsTransaction(db) && env.IsTest(e):
		a.TestTx = true
	}

	return a, nil
}

// Create creates a new user and emits event after commit
func (a aggregate) Create(ctx context.Context, u *user.Entity) (*user.Entity, error) {
	tx := a.begin()
	defer a.rollback(tx)

	res, err := repo.Create(u, tx)
	if err != nil {
		return nil, err
	}

	eventID, err := a.saveEvent(ctx, tx, res.ID, event.UserCreated, res)
	if err != nil {
		return nil, err
	}

	if err := a.commit(tx); err != nil {
		return nil, err
	}

	_ = a.publisher.Emit(ctx, eventID, event.UserCreated, res)
	return res, nil
}

// Update only updates nickname and emits event
func (a aggregate) Update(ctx context.Context, u *user.Entity) (*user.Entity, error) {
	tx := a.begin()
	defer a.rollback(tx)

	existing, err := repo.GetUserForUpdate(u.ID, tx)
	if err != nil {
		return nil, err
	}
	existing.Nickname = u.Nickname

	updated, err := repo.Update(existing, tx)
	if err != nil {
		return nil, err
	}

	eventID, err := a.saveEvent(ctx, tx, updated.ID, event.UserUpdated, updated)
	if err != nil {
		return nil, err
	}

	if err := a.commit(tx); err != nil {
		return nil, err
	}

	_ = a.publisher.Emit(ctx, eventID, event.UserUpdated, updated)
	return updated, nil
}

// Delete performs a soft delete and emits event
func (a aggregate) Delete(ctx context.Context, id uuid.UUID) error {
	tx := a.begin()
	defer a.rollback(tx)

	eventID, err := a.saveEvent(ctx, tx, id, event.UserSoftDeleted, map[string]string{
		"user_id": id.String(),
	})
	if err != nil {
		return err
	}

	if err := repo.Delete(id, tx); err != nil {
		return err
	}

	if err := a.commit(tx); err != nil {
		return err
	}

	_ = a.publisher.Emit(ctx, eventID, event.UserSoftDeleted, map[string]string{"user_id": id.String()})
	return nil
}

// Find returns a list of users with pagination and country filter
func (a aggregate) Find(_ context.Context, country string, page, limit int) ([]user.Entity, error) {
	return repo.Find(a.DB, country, page, limit)
}

func (a aggregate) begin() *gorm.DB {
	if a.TestTx {
		return a.DB
	}
	return a.DB.Begin()
}

func (a aggregate) commit(tx *gorm.DB) error {
	if !a.TestTx {
		return tx.Commit().Error
	}
	return nil
}

func (a aggregate) rollback(tx *gorm.DB) {
	if !a.TestTx {
		tx.Rollback()
	}
}

func (a *aggregate) saveEvent(ctx context.Context, tx *gorm.DB, userID uuid.UUID, eventType string, payload interface{}) (uuid.UUID, error) {
	if m, ok := payload.(map[string]string); ok {
		if traceID, ok := ctx.Value("trace_id").(string); ok {
			m["trace_id"] = traceID
		}
		payload = m
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return uuid.Nil, err
	}

	eventID := uuid.New()
	event := event.User{
		ID:        eventID,
		UserID:    userID,
		EventType: eventType,
		Payload:   data,
	}
	return eventID, tx.Create(&event).Error
}
