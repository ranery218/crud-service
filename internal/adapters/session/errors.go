package session

import "errors"

var (
	ErrInvalidTtl      = errors.New("ttl can not be 0")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)
