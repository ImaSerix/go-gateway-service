package proxy

import "errors"

var (
	ErrInvalidMethod     = errors.New("invalid method")
	ErrUnsupportedScheme = errors.New("unsupported scheme")
	ErrEmptyHost         = errors.New("empty host")
)
