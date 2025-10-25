package id_gen

import (
	"context"

	"github.com/google/uuid"
)

type DefaultIDGen struct{}

func (g *DefaultIDGen) NewID(ctx context.Context) (string, error) {
	return uuid.NewString(), nil
}
