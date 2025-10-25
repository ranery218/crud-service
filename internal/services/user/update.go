package user

import (
	"context"
	"crud/internal/domain/entities"

	"github.com/samber/mo"
)

type UpdateRequest struct {
	ID       string
	Username string
	Email    string
}

type UpdateResponse struct {
	User entities.User
}

type UpdateRepository interface {
	Update(context.Context, entities.UserAttrs, entities.UserFilterAttrs, *entities.User) error
}

type UpdateService struct {
	Repo UpdateRepository
}

func NewUpdateService(repo UpdateRepository) *UpdateService {
	return &UpdateService{
		Repo: repo,
	}
}

func (s *UpdateService) Update(ctx context.Context, req UpdateRequest) (UpdateResponse, error) {
	id := req.ID
	username := req.Username
	email := req.Email

	var updatedUser entities.User

	err := s.Repo.Update(ctx, entities.UserAttrs{
		Username: username,
		Email:    email,
	}, entities.UserFilterAttrs{
		ID: mo.Some(id),
	}, &updatedUser)
	if err != nil {
		return UpdateResponse{}, err
	}

	return UpdateResponse{
		User: updatedUser,
	}, nil
}
