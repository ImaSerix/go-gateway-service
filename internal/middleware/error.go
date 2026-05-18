package middleware

import "errors"

var (
	ErrNilConfig   = errors.New("nil config")
	ErrInvalidType = errors.New("invalid type")
)
