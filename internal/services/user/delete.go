package user

import (
	"context"
	"crud/internal/domain/entities"

	"github.com/samber/mo"
)

type DeleteRequest struct {
	ID string
}

type DeleteResponse struct {
	Success bool
}

type DeleteRepository interface {
	Delete(context.Context, entities.UserFilterAttrs) error
}

type DeleteService struct {
	Repo DeleteRepository
}

func NewDeleteService(repo DeleteRepository) *DeleteService {
	return &DeleteService{
		Repo: repo,
	}
}

func (s *DeleteService) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	err := s.Repo.Delete(ctx, entities.UserFilterAttrs{ID: mo.Some(req.ID)})
	if err != nil {
		return DeleteResponse{Success: false}, err
	}
	return DeleteResponse{Success: true}, nil
}
