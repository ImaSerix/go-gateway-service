package middleware

import "errors"

var (
	ErrEmptyAllowedOrigins        = errors.New("empty allowed origins")
	ErrInvalidAllowedMethod       = errors.New("allowed method")
	ErrInvalidLimit               = errors.New("invalid limit")
	ErrInvalidWindow              = errors.New("invalid window")
	ErrInvalidDuration            = errors.New("invalid duration")
	ErrEmptyHeaders               = errors.New("empty headers")
	ErrUnregisteredMiddlewareType = errors.New("unregisteredMiddlewareType")
)
