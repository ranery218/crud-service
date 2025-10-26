package id_gen

import (
	"github.com/google/uuid"
)

type DefaultIDGen struct{}

func NewDefaultIDGen() *DefaultIDGen {
	return &DefaultIDGen{}
}

func (g *DefaultIDGen) NewID() (string, error) {
	return uuid.NewString(), nil
}
