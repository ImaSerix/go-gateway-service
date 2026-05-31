package resolver

import (
	"net/http"
)

type Resolver interface {
	Resolve(r *http.Request, key string) (any, bool)
}
