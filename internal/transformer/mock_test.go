package transformer_test

import (
	"net/http"
	"strings"
)

type resolverMock struct {
	values        map[string]any
	forHeaderTest bool
}

func (rm *resolverMock) Resolve(r *http.Request, key string) (any, bool) {

	k := key

	if !rm.forHeaderTest {

		trim := strings.TrimPrefix(key, "{")
		trim = strings.TrimSuffix(trim, "}")

		split := strings.Split(trim, ":")

		if len(split) != 2 {
			return nil, false
		}

		k = split[1]
	}

	v, ok := rm.values[k]
	if !ok {
		return nil, false
	}
	return v, true
}
