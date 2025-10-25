package user

import (
	"context"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionStore interface {
	Create(ctx context.Context, userID string) (Session, error)
	Get(ctx context.Context, sessionID string) (Session, error)
	Delete(ctx context.Context, sessionID string) error
}

