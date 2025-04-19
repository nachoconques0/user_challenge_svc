package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
)

func TestEntity_Valid(t *testing.T) {
	validUser := &entity.User{
		ID:        uuid.New(),
		FirstName: "Nacho",
		LastName:  "Calcagno",
		Nickname:  "bandido",
		Password:  "supersecure",
		Email:     "nacho@test.com",
		Country:   "VE",
	}

	t.Run("should be valid", func(t *testing.T) {
		err := validUser.Valid()
		assert.NoError(t, err)
	})

	t.Run("should fail with empty first name", func(t *testing.T) {
		user := *validUser
		user.FirstName = ""
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrEmptyField)
	})

	t.Run("should fail with empty last name", func(t *testing.T) {
		user := *validUser
		user.LastName = "   "
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrEmptyField)
	})

	t.Run("should fail with empty nickname", func(t *testing.T) {
		user := *validUser
		user.Nickname = ""
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrEmptyField)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		user := *validUser
		user.Password = "123"
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrWeakPassword)
	})

	t.Run("should fail with empty email", func(t *testing.T) {
		user := *validUser
		user.Email = ""
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrEmptyField)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		user := *validUser
		user.Email = "invalid-email"
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrInvalidEmail)
	})

	t.Run("should fail with empty country", func(t *testing.T) {
		user := *validUser
		user.Country = "  "
		err := user.Valid()
		assert.ErrorIs(t, err, entity.ErrInvalidCountry)
	})
}

func TestEntity_HashPassword(t *testing.T) {
	user := &entity.User{
		Password: "supersecure",
	}

	err := user.HashPassword(user.Password)
	assert.NoError(t, err)
	assert.NotEqual(t, "supersecure", user.Password)
	assert.Greater(t, len(user.Password), 0)
}
