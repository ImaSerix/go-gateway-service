package checks

import "errors"

var (
	ErrNilConfig             = errors.New("nil config")
	ErrNilRequest            = errors.New("nil request")
	ErrInvalidMethod         = errors.New("ivalid method")
	ErrInvalidConfig         = errors.New("invalid config")
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrInvalidURL            = errors.New("invalid url")
	ErrEmptyHost             = errors.New("empty host")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrEmptyURL              = errors.New("empty url")
	ErrNilClient             = errors.New("nil client")
	ErrEmptyHeaders          = errors.New("empty headers")
	ErrNoHeader              = errors.New("no header")
	ErrEmptyQueries          = errors.New("empty queries")
	ErrNoQueryParam          = errors.New("no query param")
	ErrEmptyIP               = errors.New("ErrEmptyIP")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidWindow         = errors.New("invalid window")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrTooManyRequests       = errors.New("too many requests")
	ErrEmptyInjectContext    = errors.New("empty inject context")
	ErrInvalidDuration       = errors.New("invalid duration")
	ErrInvalidExpectedStatus = errors.New("invalid expected status")
)
