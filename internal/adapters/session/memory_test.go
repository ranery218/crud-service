package session

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	sessionStore, err := NewMemoryStore(30*time.Minute, nil)
	if err != nil {
		t.Fatalf("failed to create MemoryStore: %v", err)
	}

	ctx := context.Background()

	createdSession, err := sessionStore.Create(ctx, "1")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	retrievedSession, err := sessionStore.Get(ctx, createdSession.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if createdSession != retrievedSession {
		t.Fatalf("retrieved session does not match created session")
	}

	err = sessionStore.Delete(ctx, createdSession.ID)
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	_, err = sessionStore.Get(ctx, createdSession.ID)
	if err == nil {
		t.Fatalf("expected error when getting deleted session, got nil")
	}
	if err != ErrSessionNotFound {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
}

func TestMemoryStore_expiredTtl(t *testing.T) {
	sessionStore, err := NewMemoryStore(20*time.Millisecond, nil)
	if err != nil {
		t.Fatalf("failed to create MemoryStore: %v", err)
	}

	ctx := context.Background()

	createdSession, err := sessionStore.Create(ctx, "1")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	time.Sleep(40 * time.Millisecond)
	_, err = sessionStore.Get(ctx, createdSession.ID)
	if err == nil {
		t.Fatalf("expected error when getting expired session, got nil")
	}
	if err != ErrSessionExpired {
		t.Fatalf("expected ErrSessionExpired, got: %v", err)
	}
}

func TestMemoryStore_contextCanceled(t *testing.T) {
	sessionStore, err := NewMemoryStore(30*time.Minute, nil)
	if err != nil {
		t.Fatalf("failed to create MemoryStore: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = sessionStore.Create(ctx, "1")
	if err == nil {
		t.Fatalf("expected error when creating session with canceled context, got nil")
	}
}

func TestMemoryStore_NoSession(t *testing.T) {
	sessionStore, err := NewMemoryStore(30*time.Minute, nil)
	if err != nil {
		t.Fatalf("failed to create MemoryStore: %v", err)
	}

	ctx := context.Background()

	_, err = sessionStore.Get(ctx, "nonexistent-session-id")
	if err == nil {
		t.Fatalf("expected error when getting nonexistent session, got nil")
	}
	if err != ErrSessionNotFound {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
}