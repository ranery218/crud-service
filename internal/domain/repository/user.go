package repository

import (
	"context"
	"crud/internal/domain/entities"
)

type Repository interface {
	Create(context.Context, entities.UserAttrs, *entities.User) error
	Find(context.Context, entities.UserFilterAttrs, *[]entities.User) error
	FindOne(context.Context, entities.UserFilterAttrs, *entities.User) error
	Update(context.Context, entities.UserAttrs, entities.UserFilterAttrs, *entities.User) error
	Delete(context.Context, entities.UserFilterAttrs) error
}
