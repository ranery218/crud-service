package user

import (
	"context"
	"crud/internal/domain/entities"

	"github.com/samber/mo"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	User entities.User
	Session Session
}

type LoginRepository interface {
	FindOne(context.Context, entities.UserFilterAttrs, *entities.User) error
}

type LoginService struct {
	Repo   LoginRepository
	Hasher PasswordHasher
	SessionStore SessionStore
}

func NewLoginService(repo LoginRepository,hasher PasswordHasher,sessionStore SessionStore) *LoginService {
	return &LoginService{
		Repo: repo,
		Hasher: hasher,
		SessionStore: sessionStore,
	}
}

func (s *LoginService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	email := req.Email
	password := req.Password
	if email == "" {
		return LoginResponse{}, ErrEmailRequired
	}

	if password == "" {
		return LoginResponse{}, ErrPasswordRequired
	}

	email = NormalizeEmail(email)

	if !ValidateEmail(email) {
		return LoginResponse{}, ErrEmailIncorrect
	}

	var user entities.User
	err := s.Repo.FindOne(ctx, entities.UserFilterAttrs{Email: mo.Some(email)}, &user)

	if err == ErrUserNotFound {
		return LoginResponse{}, ErrUserNotFound
	}
	if err != nil {
		return LoginResponse{}, err
	}

	err = s.Hasher.Compare(ctx, user.HashedPassword, password)
	if err == ErrPasswordIncorrect {
		return LoginResponse{}, ErrPasswordIncorrect
	}
	if err != nil {
		return LoginResponse{}, err
	}

	session, err := s.SessionStore.Create(ctx, user.ID)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{User: user, Session: session}, nil
}

