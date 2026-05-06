package checks

import "errors"

var (
	ErrNilConfig       = errors.New("nil config")
	ErrNilRequest      = errors.New("nil request")
	ErrInvalidMethod   = errors.New("ivalid method")
	ErrInvalidConfig   = errors.New("invalid config")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrInvalidURL      = errors.New("invalid url")
	ErrEmptyHost       = errors.New("empty host")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrEmptyURL        = errors.New("empty url")
	ErrNilClient       = errors.New("nil client")
)
