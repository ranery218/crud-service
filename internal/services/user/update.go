package user

import (
	"context"
	"crud/internal/domain/entities"

	"github.com/samber/mo"
)

type UpdateRequest struct {
	ID       string
	Username mo.Option[string]
	Email    mo.Option[string]
	Password mo.Option[string]
}

type UpdateResponse struct {
	User entities.User
}

type UpdateRepository interface {
	Update(context.Context, entities.UserUpdateAttrs, entities.UserFilterAttrs, *entities.User) error
}
type UpdateService struct {
	Repo   UpdateRepository
	Hasher PasswordHasher
}

func NewUpdateService(repo UpdateRepository, hasher PasswordHasher) *UpdateService {
	return &UpdateService{
		Repo:   repo,
		Hasher: hasher,
	}
}

func (s *UpdateService) Update(ctx context.Context, req UpdateRequest) (UpdateResponse, error) {
	id := req.ID
	username := req.Username
	email := req.Email
	var hashedPassword mo.Option[string]

	password, ok := req.Password.Get()
	if ok {
		if !ValidatePassword(password) {
			return UpdateResponse{}, ErrPasswordIncorrect
		}
		hash, err := s.Hasher.Hash(ctx, password)
		if err != nil {
			return UpdateResponse{}, err
		}
		hashedPassword = mo.Some(hash)
	}

	var updatedUser entities.User

	err := s.Repo.Update(ctx, entities.UserUpdateAttrs{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
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
