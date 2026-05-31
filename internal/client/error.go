package client

import "errors"

var (
	ErrInvalidMethod = errors.New("invalid method")
	ErrInvalidHost   = errors.New("invalid host")
)
