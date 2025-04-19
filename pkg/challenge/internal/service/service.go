package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
)

type service struct {
	userAggregate user.Aggregate
}

type Service interface {
	Create(ctx context.Context, input *model.CreateUserInput) (*model.UserOutput, error)
	Find(ctx context.Context, country string, page, limit int) ([]model.UserOutput, error)
	Update(ctx context.Context, id uuid.UUID, nickname string) (*model.UserOutput, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func New(userAgg user.Aggregate) Service {
	return service{userAggregate: userAgg}
}

func (s service) Create(ctx context.Context, input *model.CreateUserInput) (*model.UserOutput, error) {
	userEntity := MapCreateInputToEntity(input)

	created, err := s.userAggregate.Create(ctx, userEntity)
	if err != nil {
		return nil, err
	}

	return MapEntityToOutput(created), nil
}

func (s service) Find(ctx context.Context, country string, page, limit int) ([]model.UserOutput, error) {
	users, err := s.userAggregate.Find(ctx, country, page, limit)
	if err != nil {
		return nil, err
	}

	mappedUsers := make([]model.UserOutput, 0, len(users))
	for _, u := range users {
		mappedUsers = append(mappedUsers, *MapEntityToOutput(&u))
	}

	return mappedUsers, nil
}

func (s service) Update(ctx context.Context, id uuid.UUID, nickname string) (*model.UserOutput, error) {
	updated, err := s.userAggregate.Update(ctx, &entity.User{
		ID:       id,
		Nickname: nickname,
	})
	if err != nil {
		return nil, err
	}

	return MapEntityToOutput(updated), nil
}

func (s service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userAggregate.Delete(ctx, id)
}

func MapCreateInputToEntity(in *model.CreateUserInput) *entity.User {
	return &entity.User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Nickname:  in.Nickname,
		Password:  in.Password,
		Email:     in.Email,
		Country:   in.Country,
	}
}

func MapEntityToOutput(u *entity.User) *model.UserOutput {
	return &model.UserOutput{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Country:   u.Country,
	}
}
