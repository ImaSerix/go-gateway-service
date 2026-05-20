package transformer

import "errors"

var (
	ErrNilRequest             = errors.New("nil request")
	ErrInvalidKey             = errors.New("invalid key")
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrInvalidPlaceholder     = errors.New("invalid placeholder")
	ErrInvalidListTempalate   = errors.New("invalid list template")
)
