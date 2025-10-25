package user

import "context"

type Deleterequest struct {
	ID string
}

type DeleteResponse struct {
	Success bool
}

type DeleteRepository interface {
	Delete(context.Context, string) error
}

type DeleteService struct {
	Repo DeleteRepository
}

func NewDeleteService(repo DeleteRepository) *DeleteService {
	return &DeleteService{
		Repo: repo,
	}
}

func (s *DeleteService) Delete(ctx context.Context, req Deleterequest) (DeleteResponse, error) {
	err := s.Repo.Delete(ctx, req.ID)
	if err != nil {
		return DeleteResponse{Success: false}, err
	}

	return DeleteResponse{Success: true}, nil
}
