package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	controller "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/grpc/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/mocks"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
	userProto "github.com/nachoconques0/user_challenge_svc/pkg/challenge/proto/user"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	c := controller.NewController(mockSvc)

	t.Run("should create user succes", func(t *testing.T) {
		req := &userProto.CreateUserRequest{
			FirstName: "Alice",
			LastName:  "Bob",
			Nickname:  "AB123",
			Password:  "supersecurepassword",
			Email:     "alice@bob.com",
			Country:   "UK",
		}

		expected := &model.UserOutput{
			ID:        "uuid-123",
			FirstName: "Alice",
			LastName:  "Bob",
			Nickname:  "AB123",
			Email:     "alice@bob.com",
			Country:   "UK",
		}

		mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expected, nil)

		res, err := c.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", res.FirstName)
	})

	t.Run("should return error when missing fields", func(t *testing.T) {
		req := &userProto.CreateUserRequest{}

		_, err := c.CreateUser(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	c := controller.NewController(mockSvc)

	t.Run("should update nickname", func(t *testing.T) {
		req := &userProto.UpdateUserRequest{
			Id:       "d4e3e4ea-6a0b-4c2e-9e5c-cd6fdf2de771",
			Nickname: "newcsgoplayer",
		}

		mockSvc.EXPECT().
			Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&model.UserOutput{Nickname: "newcsgoplayer"}, nil)

		res, err := c.UpdateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "newcsgoplayer", res.Nickname)
	})

	t.Run("should fail on invalid UUID", func(t *testing.T) {
		req := &userProto.UpdateUserRequest{
			Id:       "not-a-uuid",
			Nickname: "wannaplaycsgo",
		}
		_, err := c.UpdateUser(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	c := controller.NewController(mockSvc)

	t.Run("should delete user", func(t *testing.T) {
		req := &userProto.DeleteUserRequest{Id: "d4e3e4ea-6a0b-4c2e-9e5c-cd6fdf2de771"}

		mockSvc.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Return(nil)

		res, err := c.DeleteUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("should fail on invalid ID", func(t *testing.T) {
		req := &userProto.DeleteUserRequest{Id: "invalid-id"}

		_, err := c.DeleteUser(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestFindUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	c := controller.NewController(mockSvc)

	t.Run("should find users", func(t *testing.T) {
		req := &userProto.FindUsersRequest{
			Country: "UK",
			Page:    1,
			Limit:   2,
		}

		expected := []model.UserOutput{
			{ID: "1", FirstName: "nachooo"},
			{ID: "2", FirstName: "faceit"},
		}

		mockSvc.EXPECT().
			Find(gomock.Any(), "UK", 1, 2).
			Return(expected, nil)

		res, err := c.FindUsers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, res.Users, 2)
	})

	t.Run("should handle service error", func(t *testing.T) {
		req := &userProto.FindUsersRequest{
			Country: "UK",
			Page:    1,
			Limit:   2,
		}

		mockSvc.EXPECT().
			Find(gomock.Any(), "UK", 1, 2).
			Return(nil, errors.New("errtest"))

		_, err := c.FindUsers(context.Background(), req)
		assert.Error(t, err)
	})
}
