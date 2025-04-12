package errors

import "errors"

var (
	ErrSessionNotExists   = errors.New("such session isn't exist")
	ErrSessionRevoked     = errors.New("session is revoked")
	ErrSessionExpired     = errors.New("session is expired")
	ErrInternal           = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccessDenied       = errors.New("access denied")
)
