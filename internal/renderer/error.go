package renderer

import "errors"

var (
	ErrInvalidTemplate = errors.New("invalid template")
	ErrFailedResolve   = errors.New("failed resolve")
)
