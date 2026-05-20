package resolver

import "net/http"

type HeaderResolver struct{}

func NewHeaderResolver() *HeaderResolver {
	return &HeaderResolver{}
}

func (cr *HeaderResolver) Resolve(r *http.Request, key string) (any, bool) {

	_, n, err := getSourceAndKey(key)
	if err != nil {
		return nil, false
	}

	if v := r.Header.Get(n); v != "" {
		return v, true
	}

	return nil, false
}
