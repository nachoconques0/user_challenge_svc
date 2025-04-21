package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	userAgg "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
)

type service struct {
	userAggregate userAgg.Aggregate
}

type Service interface {
	Create(ctx context.Context, input *model.CreateUserInput) (*model.UserOutput, error)
	Find(ctx context.Context, country string, page, limit int) ([]model.UserOutput, error)
	Update(ctx context.Context, id uuid.UUID, nickname string) (*model.UserOutput, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// New returns a new User service
func New(userAgg userAgg.Aggregate) Service {
	return service{userAggregate: userAgg}
}

// Create creates a new user and emits event after commit
func (s service) Create(ctx context.Context, input *model.CreateUserInput) (*model.UserOutput, error) {
	userEntity := mapCreateInputToEntity(input)

	created, err := s.userAggregate.Create(ctx, userEntity)
	if err != nil {
		log.Error().Err(err).Str("userService", "Create").Msg("could not create user")
		return nil, err
	}

	return mapEntityToOutput(created), nil
}

// Find returns a list of users with pagination and country filter
func (s service) Find(ctx context.Context, country string, page, limit int) ([]model.UserOutput, error) {
	users, err := s.userAggregate.Find(ctx, country, page, limit)
	if err != nil {
		log.Error().Err(err).Str("userService", "Find").Msg("could not find users")
		return nil, err
	}

	mappedUsers := make([]model.UserOutput, 0, len(users))
	for _, u := range users {
		mappedUsers = append(mappedUsers, *mapEntityToOutput(&u))
	}

	return mappedUsers, nil
}

// Update only updates nickname and emits event
func (s service) Update(ctx context.Context, id uuid.UUID, nickname string) (*model.UserOutput, error) {
	updated, err := s.userAggregate.Update(ctx, &user.Entity{
		ID:       id,
		Nickname: nickname,
	})
	if err != nil {
		log.Error().Err(err).Str("userService", "Update").Msg("could not update user")
		return nil, err
	}

	return mapEntityToOutput(updated), nil
}

// Delete performs a soft delete and emits event
func (s service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.userAggregate.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("userService", "Delete").Msg("could not soft delete user")
		return err
	}
	return nil
}

func mapCreateInputToEntity(in *model.CreateUserInput) *user.Entity {
	return &user.Entity{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Nickname:  in.Nickname,
		Password:  in.Password,
		Email:     in.Email,
		Country:   in.Country,
	}
}

func mapEntityToOutput(u *user.Entity) *model.UserOutput {
	return &model.UserOutput{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Country:   u.Country,
	}
}
