package endpoint

import "errors"

var (
	ErrEmptyPath         = errors.New("empty path")
	ErrPathHasSpaces     = errors.New("path has spaces")
	ErrUnsupportedScheme = errors.New("unsupported scheme")
	ErrEmptyHost         = errors.New("empty host")
	ErrInvalidMethod     = errors.New("invalid method")
	ErrInvalidConfig     = errors.New("invalid config")
)
