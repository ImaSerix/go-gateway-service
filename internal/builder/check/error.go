package check

import "errors"

var (
	ErrUnregisteredCheckName = errors.New("unregistered check name")
)
