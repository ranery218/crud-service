package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailTaken        = errors.New("email already in use")
	ErrUsernameTaken     = errors.New("username already in use")
	ErrUsernameRequired  = errors.New("username is required")
	ErrEmailRequired     = errors.New("email is required")
	ErrPasswordRequired  = errors.New("password is required")
	ErrEmailIncorrect    = errors.New("incorrect email")
	ErrPasswordIncorrect = errors.New("incorrect password")
	ErrSessionNotFound   = errors.New("session not found")
	ErrSessionExpired    = errors.New("session is expired")
)
