package helpers

import "errors"

var (
	ErrBodyEmpty = errors.New("body is empty")
	ErrMalformedJSON = errors.New("json is malformed")
	ErrInvalidType = errors.New("invalid type")
)