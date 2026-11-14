package renderer

import "errors"

var (
	ErrInvalidTemplate = errors.New("invalid template")
	ErrFailedResolve   = errors.New("failed resolve")
	ErrDotSegment      = errors.New("value makes a dot segment in the path")
)
