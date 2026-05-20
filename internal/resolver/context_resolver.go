package resolver

import "net/http"

type ContextResolver struct{}

func NewContextResolver() *ContextResolver {
	return &ContextResolver{}
}

func (cr *ContextResolver) Resolve(r *http.Request, key string) (any, bool) {

	_, n, err := getSourceAndKey(key)
	if err != nil {
		return nil, false
	}

	if v := r.Context().Value(n); v != nil {
		return v, true
	}

	return nil, false
}
