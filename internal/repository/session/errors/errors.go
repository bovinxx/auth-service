package errors

import "errors"

var (
	ErrSessionNotExists = errors.New("such session isn't exist")
)
