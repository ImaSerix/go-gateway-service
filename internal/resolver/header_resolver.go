package resolver

import "net/http"

type HeaderResolver struct{}

func NewHeaderResolver() *HeaderResolver {
	return &HeaderResolver{}
}

func (cr *HeaderResolver) Resolve(r *http.Request, key string) (any, bool) {

	if v := r.Header.Get(key); v != "" {
		return v, true
	}

	return nil, false
}
