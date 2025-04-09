package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("such user already exists")
	ErrUserNotExists      = errors.New("such user doesn't exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
