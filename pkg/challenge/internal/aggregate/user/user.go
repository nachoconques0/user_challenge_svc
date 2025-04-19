package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/google/uuid"
	dbInstance "github.com/nachoconques0/user_challenge_svc/pkg/challenge/db"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/env"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/repo"
)

type aggregate struct {
	DB     *gorm.DB
	TestTx bool
}

type Aggregate interface {
	Create(ctx context.Context, u *entity.User) (*entity.User, error)
	Find(ctx context.Context, country string, page, limit int) ([]entity.User, error)
	Update(ctx context.Context, u *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

const (
	ErrMissingDB      = "Aggregate is missing DB connection"
	ErrMissingTestEnv = "DB connection can only be a TX when ENV == env.Test"
)

func New(db *gorm.DB, e string) (aggregate, error) {
	// Init the aggregate
	a := aggregate{
		DB: db,
	}

	// Make sure that the Aggregate has all the options needed
	var err error
	if db == nil {
		err = errors.New(ErrMissingDB)
	} else if dbInstance.IsTransaction(db) {
		// We don't want DB transactions to be passed unless the
		if env.IsTest(e) {
			// running environment is test
			a.TestTx = true
		} else {
			err = errors.New(ErrMissingTestEnv)
		}
	}

	return a, err
}

// Create creates a new user in the DB
func (a aggregate) Create(_ context.Context, u *entity.User) (*entity.User, error) {
	tx := a.begin()
	defer a.rollback(tx)

	res, err := repo.Create(u, tx)
	if err != nil {
		return nil, err
	}

	err = a.commit(tx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Find returns a list of users. It can be paginated and filtered by user country
func (a aggregate) Find(_ context.Context, country string, page, limit int) ([]entity.User, error) {
	return repo.Find(a.begin(), country, page, limit)
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

func (a aggregate) Update(_ context.Context, u *entity.User) (*entity.User, error) {
	tx := a.begin()
	defer a.rollback(tx)

	userForUpdate, err := repo.GetUserForUpdate(u.ID, tx)
	if err != nil {
		return nil, err
	}

	userForUpdate.Nickname = u.Nickname
	res, err := repo.Update(userForUpdate, tx)
	if err != nil {
		return nil, err
	}

	err = a.commit(tx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a aggregate) Delete(_ context.Context, id uuid.UUID) error {
	tx := a.begin()
	defer a.rollback(tx)

	if err := repo.Delete(id, tx); err != nil {
		return err
	}

	return a.commit(tx)
}
