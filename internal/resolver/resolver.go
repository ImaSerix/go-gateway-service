package resolver

import (
	"net/http"
	"strings"
)

type Resolver interface {
	Resolve(r *http.Request, key string) (any, bool)
}

func getSourceAndKey(key string) (source string, name string, err error) {

	trim := strings.TrimSpace(key)

	split := strings.Split(trim, ".")

	if len(split) != 2 {
		return "", "", ErrInvalidKeyFormat
	}

	if split[0] == "" {
		return "", "", ErrEmptyKeySource
	}

	if split[1] == "" {
		return "", "", ErrEmptyKeyName
	}

	return split[0], split[1], nil
}
