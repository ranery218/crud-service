package user

import (
	"context"
	"crud/internal/domain/entities"
	"errors"
	"testing"
	"time"
)

type loginRepoStub struct {
	user entities.User
	err  error
}

func (r *loginRepoStub) FindOne(ctx context.Context, attrs entities.UserFilterAttrs, ent *entities.User) error {
	if r.err != nil {
		return r.err
	}

	email, ok := attrs.Email.Get()
	if !ok || email == "" {
		return ErrEmailRequired
	}

	*ent = r.user
	return nil
}

type hasherStub struct {
	compareErr error
}

func (h *hasherStub) Hash(ctx context.Context, plaintext string) (string, error) {
	return "hashed", nil
}

func (h *hasherStub) Compare(ctx context.Context, hash string, plaintext string) error {
	return h.compareErr
}

type sessionStoreStub struct {
	session   Session
	createErr error
}

func (s *sessionStoreStub) Create(ctx context.Context, userID string) (Session, error) {
	if s.createErr != nil {
		return Session{}, s.createErr
	}

	if s.session.ID == "" {
		s.session = Session{
			ID:        "session-1",
			UserID:    userID,
			ExpiresAt: time.Time{},
		}
	}
	return s.session, nil
}

func (s *sessionStoreStub) Get(ctx context.Context, sessionID string) (Session, error) {
	return Session{}, nil
}

func (s *sessionStoreStub) Delete(ctx context.Context, sessionID string) error {
	return nil
}

func TestLogin_Success(t *testing.T) {
	repo := &loginRepoStub{
		user: entities.User{
			ID:             "1",
			Email:          "islam@gmail.com",
			Username:       "islam",
			HashedPassword: "hashed",
		},
	}
	hasher := &hasherStub{compareErr: nil}
	store := &sessionStoreStub{}
	loginService := NewLoginService(repo, hasher, store)

	ctx := context.Background()
	request := LoginRequest{Email: "islam@gmail.com", Password: "secret"}

	response, err := loginService.Login(ctx, request)
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	if response.Session.ID == "" || response.Session.UserID != repo.user.ID {
		t.Fatalf("unexpected session: %+v", response.Session)
	}

	if response.User != repo.user {
		t.Fatalf("unexpected user: %+v", response.User)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := &loginRepoStub{
		user: entities.User{
			ID:             "1",
			Email:          "islam@gmail.com",
			Username:       "islam",
			HashedPassword: "hashed",
		},
	}
	hasher := &hasherStub{compareErr: ErrPasswordIncorrect}
	store := &sessionStoreStub{}
	loginService := NewLoginService(repo, hasher, store)

	_, err := loginService.Login(context.Background(), LoginRequest{
		Email:    "islam@gmail.com",
		Password: "wrong",
	})
	if !errors.Is(err, ErrPasswordIncorrect) {
		t.Fatalf("expected ErrPasswordIncorrect, got: %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &loginRepoStub{
		err: ErrUserNotFound,
	}
	hasher := &hasherStub{}
	store := &sessionStoreStub{}
	loginService := NewLoginService(repo, hasher, store)

	_, err := loginService.Login(context.Background(), LoginRequest{
		Email:    "missing@example.com",
		Password: "secret",
	})
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestLogin_SessionCreationFailure(t *testing.T) {
	repo := &loginRepoStub{
		user: entities.User{
			ID:             "1",
			Email:          "islam@gmail.com",
			Username:       "islam",
			HashedPassword: "hashed",
		},
	}
	hasher := &hasherStub{}
	createErr := errors.New("session create failed")
	store := &sessionStoreStub{createErr: createErr}
	loginService := NewLoginService(repo, hasher, store)

	_, err := loginService.Login(context.Background(), LoginRequest{
		Email:    "islam@gmail.com",
		Password: "secret",
	})
	if !errors.Is(err, createErr) {
		t.Fatalf("expected session create error, got: %v", err)
	}
}
