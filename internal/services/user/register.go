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
	Create(ctx context.Context, attrs entities.UserAttrs, ent *entities.User) error
	FindOne(ctx context.Context, filterAttrs entities.UserFilterAttrs, ent *entities.User) error
}

type RegisterService struct {
	Repo   RegisterRepository
	Hasher PasswordHasher
	IdGen  IDGen
}

func NewRegisterService(repo RegisterRepository, hasher PasswordHasher, idGen IDGen) *RegisterService {
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

	err := s.Repo.FindOne(ctx, entities.UserFilterAttrs{
		Email: mo.Some(email),
	}, &entities.User{})
	if err != ErrUserNotFound {
		if err == nil {
			return RegisterResponse{}, ErrEmailTaken
		}
		return RegisterResponse{}, err
	}

	err = s.Repo.FindOne(ctx, entities.UserFilterAttrs{
		Username: mo.Some(username),
	}, &entities.User{})
	if err != ErrUserNotFound {
		if err == nil {
			return RegisterResponse{}, ErrUsernameTaken
		}
		return RegisterResponse{}, err
	}

	id, err := s.IdGen.NewID()
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
