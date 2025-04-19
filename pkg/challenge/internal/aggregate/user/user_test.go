package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/helpers"
	agg "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
)

func TestAggregate_UserLifecycle(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := helpers.NewTestDB()
	assert.NoError(t, err)
	defer teardown()

	aggregate, err := agg.New(db, "test")
	assert.NoError(t, err)

	// we create an user
	user := &entity.User{
		FirstName: "Test",
		LastName:  "User",
		Nickname:  "bandido",
		Password:  "securepass123",
		Email:     "test@bandido.com",
		Country:   "VE",
	}

	created, err := aggregate.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "bandido", created.Nickname)

	// we find users
	users, err := aggregate.Find(ctx, "VE", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, created.ID, users[0].ID)

	// we update the user
	created.Nickname = "bandidoooooo"
	updated, err := aggregate.Update(ctx, created)
	assert.NoError(t, err)
	assert.Equal(t, "bandidoooooo", updated.Nickname)

	// we soft delete the user
	err = aggregate.Delete(ctx, created.ID)
	assert.NoError(t, err)

	usersAfterSoftDelete, err := aggregate.Find(ctx, "", 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, usersAfterSoftDelete)
}
