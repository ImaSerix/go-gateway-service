package transformer

import "errors"

var (
	ErrInvalidConfig             = errors.New("invalid config")
	ErrUnregisteredTransformName = errors.New("unregistered transform name")
)
