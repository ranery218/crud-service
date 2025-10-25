package user

import "context"

type IDGen interface {
	NewID(context.Context) (string,error)
}