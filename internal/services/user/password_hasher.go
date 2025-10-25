package user

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, plaintext string) (string, error)
	Compare(ctx context.Context, hash string, plaintext string) error
}
