package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
	service "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service/user"
	userProto "github.com/nachoconques0/user_challenge_svc/pkg/challenge/proto/user"
)

var (
	// ErrMissingFields used user request is missing fields
	ErrMissingFields = errors.New("missing fields")
	// ErrIDnotValid used when entity ID is not valid
	ErrIDnotValid = errors.New("ID is not valid")
)

type Controller struct {
	svc service.Service
	userProto.UnimplementedUserServiceServer
}

// NewController returns a gRPC User controller
func NewController(s service.Service) *Controller {
	return &Controller{svc: s}
}

// CreateUser returns a created user
func (c *Controller) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.UserResponse, error) {
	if strings.TrimSpace(req.FirstName) == "" ||
		strings.TrimSpace(req.LastName) == "" ||
		strings.TrimSpace(req.Nickname) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Country) == "" {
		log.Error().Err(ErrMissingFields).Str("userController", "CreateUser").Msg("not valid data")
		return nil, ErrMissingFields
	}

	in := &model.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}
	user, err := c.svc.Create(ctx, in)
	if err != nil {
		log.Error().Err(ErrMissingFields).Str("userController", "CreateUser").Msg("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return mapToProto(user), nil
}

// UpdateUser updates an user nickname
func (c *Controller) UpdateUser(ctx context.Context, req *userProto.UpdateUserRequest) (*userProto.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, ErrIDnotValid
	}
	if strings.TrimSpace(req.Nickname) == "" {
		return nil, ErrMissingFields
	}
	user, err := c.svc.Update(ctx, id, req.Nickname)
	if err != nil {
		log.Error().Err(ErrMissingFields).Str("userController", "UpdateUser").Msg("failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return mapToProto(user), nil
}

// DeleteUser soft deletes an user based on its ID
func (c *Controller) DeleteUser(ctx context.Context, req *userProto.DeleteUserRequest) (*userProto.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	if err := c.svc.Delete(ctx, id); err != nil {
		log.Error().Err(ErrMissingFields).Str("userController", "DeleteUser").Msg("failed to delete user")
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return &userProto.Empty{}, nil
}

// FindUsers returns a list of users. It is paginated and also can be filtered by country
func (c *Controller) FindUsers(ctx context.Context, req *userProto.FindUsersRequest) (*userProto.UsersResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	users, err := c.svc.Find(ctx, req.Country, int(req.Page), int(req.Limit))
	if err != nil {
		log.Error().Err(ErrMissingFields).Str("userController", "FindUsers").Msg("failed to find users")
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	res := &userProto.UsersResponse{}
	for _, u := range users {
		res.Users = append(res.Users, mapToProto(&u))
	}
	return res, nil
}

func mapToProto(u *model.UserOutput) *userProto.UserResponse {
	return &userProto.UserResponse{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Country:   u.Country,
	}
}
