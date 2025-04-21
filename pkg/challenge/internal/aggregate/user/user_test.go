package user_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/helpers"
	agg "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user"
	eventUser "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user/event"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/mocks"
)

func TestUserAggregate_Create(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	assert.NoError(t, err)
	defer teardown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := mocks.NewMockPublisher(ctrl)
	aggregate, err := agg.New(db, "test", mockPublisher)
	assert.NoError(t, err)

	ctx := context.Background()
	input := &user.Entity{
		FirstName: "Nacho",
		LastName:  "Calcagno",
		Nickname:  "bandido123",
		Password:  "123123123",
		Email:     "nacho@gmail.com",
		Country:   "VE",
	}

	mockPublisher.EXPECT().Emit(gomock.Any(), gomock.Any(), eventUser.UserCreated, gomock.Any()).Return(nil)

	created, err := aggregate.Create(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, "Nacho", created.FirstName)

	var event eventUser.User
	err = db.Where("user_id = ? AND event_type = ?", created.ID, eventUser.UserCreated).First(&event).Error
	assert.NoError(t, err)

	var payload user.Entity
	err = json.Unmarshal(event.Payload, &payload)
	assert.NoError(t, err)
	assert.Equal(t, created.Email, payload.Email)
}

func TestUserAggregate_Update(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	assert.NoError(t, err)
	defer teardown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := mocks.NewMockPublisher(ctrl)
	agg, err := agg.New(db, "test", mockPublisher)
	assert.NoError(t, err)

	ctx := context.Background()
	user := &user.Entity{
		FirstName: "Juan",
		LastName:  "Perez",
		Nickname:  "jperez",
		Password:  "12345678",
		Email:     "juanp@correo.com",
		Country:   "AR",
	}

	mockPublisher.EXPECT().Emit(gomock.Any(), gomock.Any(), eventUser.UserCreated, gomock.Any()).Return(nil)
	created, err := agg.Create(ctx, user)
	assert.NoError(t, err)

	created.Nickname = "csgooo"
	mockPublisher.EXPECT().Emit(gomock.Any(), gomock.Any(), eventUser.UserUpdated, gomock.Any()).Return(nil)

	updated, err := agg.Update(ctx, created)
	assert.NoError(t, err)
	assert.Equal(t, "csgooo", updated.Nickname)

	var event eventUser.User
	err = db.Where("user_id = ? AND event_type = ?", updated.ID, eventUser.UserUpdated).First(&event).Error
	assert.NoError(t, err)
}

func TestUserAggregate_Delete(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	assert.NoError(t, err)
	defer teardown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := mocks.NewMockPublisher(ctrl)
	agg, err := agg.New(db, "test", mockPublisher)
	assert.NoError(t, err)

	ctx := context.Background()
	u := &user.Entity{
		FirstName: "Juan",
		LastName:  "calcagno",
		Nickname:  "softsoft",
		Password:  "12345678",
		Email:     "abouttodeleeee@gmail.com",
		Country:   "ES",
	}

	mockPublisher.EXPECT().Emit(gomock.Any(), gomock.Any(), eventUser.UserCreated, gomock.Any()).Return(nil)
	created, err := agg.Create(ctx, u)
	assert.NoError(t, err)

	mockPublisher.EXPECT().Emit(gomock.Any(), gomock.Any(), eventUser.UserSoftDeleted, gomock.Any()).Return(nil)

	err = agg.Delete(ctx, created.ID)
	assert.NoError(t, err)

	var event eventUser.User
	err = db.Where("user_id = ? AND event_type = ?", created.ID, eventUser.UserSoftDeleted).First(&event).Error
	assert.NoError(t, err)
}
