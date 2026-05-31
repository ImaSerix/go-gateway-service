package resolver

import "errors"

var (
	ErrInvalidKeyFormat = errors.New("invalid key format")
	ErrEmptyKeySource   = errors.New("empty key source")
	ErrEmptyKeyName     = errors.New("empty key name")
)
