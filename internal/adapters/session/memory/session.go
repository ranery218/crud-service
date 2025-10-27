package memory

import (
	"context"
	"crud/internal/services/user"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]user.Session
	ttl      time.Duration
	idGen    func() (string, error)
}

func NewMemoryStore(ttl time.Duration, idGen func() (string, error)) (*MemoryStore, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTtl
	}
	if idGen == nil {
		idGen = func() (string, error) {
			return uuid.NewString(), nil
		}
	}
	return &MemoryStore{
		sessions: make(map[string]user.Session),
		ttl:      ttl,
		idGen:    idGen,
	}, nil
}

func (s *MemoryStore) Create(ctx context.Context, userID string) (user.Session, error) {
	if err := ctx.Err(); err != nil {
		return user.Session{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := s.idGen()
	if err != nil {
		return user.Session{}, err
	}

	session := user.Session{ID: id, UserID: userID, ExpiresAt: time.Now().UTC().Add(s.ttl)}
	s.sessions[id] = session

	return session, nil
}

func (s *MemoryStore) Get(ctx context.Context, sessionID string) (user.Session, error) {
	if err := ctx.Err(); err != nil {
		return user.Session{}, err
	}

	s.mu.RLock()
	session, ok := s.sessions[sessionID]
	s.mu.RUnlock()

	if !ok {
		return user.Session{}, ErrSessionNotFound
	}

	timeNow := time.Now().UTC()
	if timeNow.After(session.ExpiresAt) {
		s.mu.Lock()
		if sess, ok := s.sessions[sessionID]; ok && timeNow.After(sess.ExpiresAt) {
			delete(s.sessions, sessionID)
		}
		s.mu.Unlock()
		return user.Session{}, ErrSessionExpired
	}

	return session, nil
}

func (s *MemoryStore) Delete(ctx context.Context, sessionID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	delete(s.sessions, sessionID)
	return nil
}
