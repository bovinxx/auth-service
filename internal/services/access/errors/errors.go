package errors

import "errors"

var (
	ErrNoMetadata       = errors.New("metadata not provided")
	ErrNoAuthHeader     = errors.New("authorization header not provided")
	ErrInvalidAuth      = errors.New("invalid authorization header format")
	ErrInvalidToken     = errors.New("access token is invalid")
	ErrAccessDenied     = errors.New("access denied")
	ErrSessionNotExists = errors.New("session not exists")
)
