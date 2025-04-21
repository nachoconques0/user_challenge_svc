package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/mocks"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
	service "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service/user"
)

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	input := &model.CreateUserInput{
		FirstName: "Nacho",
		LastName:  "Calcagno",
		Nickname:  "bandido",
		Password:  "123123123",
		Email:     "nacho@gmail.com",
		Country:   "VE",
	}

	expected := &user.Entity{
		ID:        uuid.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Nickname:  input.Nickname,
		Email:     input.Email,
		Country:   input.Country,
	}

	mockAgg.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(expected, nil)

	result, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.Nickname, result.Nickname)
}

func TestService_Create_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	input := &model.CreateUserInput{
		FirstName: "Nacho",
		LastName:  "Error",
		Nickname:  "fail",
		Password:  "12345678",
		Email:     "fail@user.com",
		Country:   "VE",
	}

	mockAgg.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	res, err := svc.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_Find_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	mockUsers := []user.Entity{
		{
			ID:        uuid.New(),
			FirstName: "A",
			LastName:  "B",
			Nickname:  "nacho",
			Email:     "a@b.com",
			Country:   "VE",
		},
	}

	mockAgg.EXPECT().
		Find(gomock.Any(), "VE", 1, 10).
		Return(mockUsers, nil)

	res, err := svc.Find(context.Background(), "VE", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "nacho", res[0].Nickname)
}

func TestService_Find_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	mockAgg.EXPECT().
		Find(gomock.Any(), "VE", 1, 10).
		Return(nil, errors.New("db failure"))

	res, err := svc.Find(context.Background(), "VE", 1, 10)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	id := uuid.New()
	nick := "cslover"

	entityUser := &user.Entity{ID: id, Nickname: nick}

	mockAgg.EXPECT().
		Update(gomock.Any(), entityUser).
		Return(entityUser, nil)

	res, err := svc.Update(context.Background(), id, nick)
	assert.NoError(t, err)
	assert.Equal(t, nick, res.Nickname)
}

func TestService_Update_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	id := uuid.New()
	testNickName := "nachin"
	user := &user.Entity{ID: id, Nickname: testNickName}

	mockAgg.EXPECT().
		Update(gomock.Any(), user).
		Return(nil, errors.New("update failed"))

	res, err := svc.Update(context.Background(), id, testNickName)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	id := uuid.New()
	mockAgg.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil)

	err := svc.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func TestService_Delete_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgg := mocks.NewMockUserAggregate(ctrl)
	svc := service.New(mockAgg)

	id := uuid.New()
	mockAgg.EXPECT().
		Delete(gomock.Any(), id).
		Return(errors.New("delete error"))

	err := svc.Delete(context.Background(), id)
	assert.Error(t, err)
}
