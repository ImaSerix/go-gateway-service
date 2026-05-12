package transformer

import "errors"

var (
	ErrNilRequest             = errors.New("nil request")
	ErrNoKeyInContext         = errors.New("no key in context")
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrInvalidPlaceholder     = errors.New("invalid placeholder")
	ErrInvalidListTempalate   = errors.New("invalid list template")
)
