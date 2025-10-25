package user

import (
	"context"
	"crud/internal/domain/entities"

	"github.com/samber/mo"
)

type GetRequest struct {
	ID string
}

type GetResponse struct {
	User entities.User
}

type GetRepository interface {
	FindOne(context.Context, entities.UserFilterAttrs, *entities.User) error
}

type GetService struct {
	Repo GetRepository
}

func NewGetService(repo GetRepository) *GetService {
	return &GetService{
		Repo: repo,
	}
}

func (s *GetService) Get(ctx context.Context, req GetRequest) (GetResponse, error) {
	var user entities.User

	err := s.Repo.FindOne(ctx, entities.UserFilterAttrs{
		ID: mo.Some(req.ID),
	}, &user)
	if err != nil {
		return GetResponse{}, err
	}

	return GetResponse{User: user}, nil
}
