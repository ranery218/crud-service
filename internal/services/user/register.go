package user

import (
	"context"
	"strings"

	"github.com/samber/mo"

	"crud/internal/domain/entities"
)

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type RegisterResponse struct {
	User entities.User
}

type RegisterRepository interface {
	Exists(ctx context.Context, filter entities.UserFilterAttrs) (bool, error)
	Create(ctx context.Context, attrs entities.UserAttrs, ent *entities.User) error
}

type IDGenerator interface {
	NewID(ctx context.Context) (string, error)
}

type RegisterService struct {
	Repo   RegisterRepository
	Hasher PasswordHasher
	IdGen  IDGenerator
}

func NewRegisterService(repo RegisterRepository, hasher PasswordHasher, idGen IDGenerator) *RegisterService {
	return &RegisterService{
		Repo:   repo,
		Hasher: hasher,
		IdGen:  idGen,
	}
}

func (s *RegisterService) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	username := strings.TrimSpace(req.Username)
	password := req.Password

	if email == "" {
		return RegisterResponse{}, ErrEmailRequired
	}

	if password == "" {
		return RegisterResponse{}, ErrPasswordRequired
	}

	if username == "" {
		return RegisterResponse{}, ErrUsernameRequired
	}

	if !ValidateEmail(email) {
		return RegisterResponse{}, ErrEmailIncorrect
	}

	if !ValidatePassword(password) {
		return RegisterResponse{}, ErrPasswordIncorrect
	}

	emailExists, err := s.Repo.Exists(ctx, entities.UserFilterAttrs{
		Email: mo.Some(email),
	})
	if err != nil {
		return RegisterResponse{}, err
	}
	if emailExists {
		return RegisterResponse{}, ErrEmailTaken
	}

	usernameExists, err := s.Repo.Exists(ctx, entities.UserFilterAttrs{
		Username: mo.Some(username),
	})
	if err != nil {
		return RegisterResponse{}, err
	}
	if usernameExists {
		return RegisterResponse{}, ErrUsernameTaken
	}

	id, err := s.IdGen.NewID(ctx)
	if err != nil {
		return RegisterResponse{}, err
	}

	hashedPassword, err := s.Hasher.Hash(ctx, password)
	if err != nil {
		return RegisterResponse{}, err
	}

	attrs := entities.UserAttrs{
		ID:             id,
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	var user entities.User
	err = s.Repo.Create(ctx, attrs, &user)
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{User: user}, nil
}
