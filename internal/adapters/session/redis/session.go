package redis

import (
	"context"
	"crud/internal/services/user"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
	idGen  func() (string, error)
}

func NewRedisStore(client *redis.Client, ttl time.Duration, idGen func() (string, error)) *RedisStore {
	return &RedisStore{
		client: client,
		ttl:    ttl,
		idGen:  idGen,
	}
}

func (s *RedisStore) Create(ctx context.Context, userID string) (user.Session, error) {
	if err := ctx.Err(); err != nil {
		return user.Session{}, err
	}
	id, err := s.idGen()
	if err != nil {
		return user.Session{}, err
	}
	session := user.Session{ID: id, UserID: userID, ExpiresAt: time.Now().UTC().Add(s.ttl)}
	payload, err := json.Marshal(session)
	if err != nil {
		return user.Session{}, err
	}
	err = s.client.Set(ctx, fmt.Sprintf("session:%s", session.ID), payload, s.ttl).Err()
	if err != nil {
		return user.Session{}, err
	}
	return session, nil
}

func (s *RedisStore) Get(ctx context.Context, sessionID string) (user.Session, error) {
	if ctx.Err() != nil {
		return user.Session{}, ctx.Err()
	}

	data, err := s.client.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return user.Session{}, user.ErrSessionNotFound
		}
		return user.Session{}, err
	}

	var session user.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return user.Session{}, err
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return user.Session{}, user.ErrSessionNotFound
	}

	return session, nil
}

func (s *RedisStore) Delete(ctx context.Context, sessionID string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := s.client.Del(ctx, fmt.Sprintf("session:%s", sessionID)).Err(); err != nil {
		if err == redis.Nil {
			return user.ErrSessionNotFound
		}
		return err
	}

	return nil
}
