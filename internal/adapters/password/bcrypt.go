package password

import (
	"context"
	"crud/internal/services/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(ctx context.Context, plaintext string) (string, error) {

	if err := ctx.Err(); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), h.cost)
	return string(hash), err

}

func (h *BcryptHasher) Compare(ctx context.Context, hash string, plaintext string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return user.ErrPasswordIncorrect
	}
	return err

}
