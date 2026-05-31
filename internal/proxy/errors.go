package proxy

import "errors"

var (
	ErrInvalidMethod     = errors.New("invalid method")
	ErrUnsupportedScheme = errors.New("unsupported scheme")
	ErrEmptyHost         = errors.New("empty host")
	ErrInvalidScheme     = errors.New("invalid scheme")
	ErrUnresolverScheme  = errors.New("unresolved scheme")
	ErrInvalidHost       = errors.New("invalid host")
)
