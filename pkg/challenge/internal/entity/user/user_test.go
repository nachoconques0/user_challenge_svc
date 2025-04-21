package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user"
	"github.com/stretchr/testify/assert"
)

func TestUserEntity_Valid(t *testing.T) {
	validUser := &user.Entity{
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
		u := *validUser
		u.FirstName = ""
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrEmptyField)
	})

	t.Run("should fail with empty last name", func(t *testing.T) {
		u := *validUser
		u.LastName = "   "
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrEmptyField)
	})

	t.Run("should fail with empty nickname", func(t *testing.T) {
		u := *validUser
		u.Nickname = ""
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrEmptyField)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		u := *validUser
		u.Password = "123"
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrWeakPassword)
	})

	t.Run("should fail with empty email", func(t *testing.T) {
		u := *validUser
		u.Email = ""
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrEmptyField)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		u := *validUser
		u.Email = "invalid-email"
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrInvalidEmail)
	})

	t.Run("should fail with empty country", func(t *testing.T) {
		u := *validUser
		u.Country = "  "
		err := u.Valid()
		assert.ErrorIs(t, err, user.ErrInvalidCountry)
	})
}

func TestEntity_HashPassword(t *testing.T) {
	user := &user.Entity{
		Password: "supersecure",
	}

	err := user.HashPassword(user.Password)
	assert.NoError(t, err)
	assert.NotEqual(t, "supersecure", user.Password)
	assert.Greater(t, len(user.Password), 0)
}
