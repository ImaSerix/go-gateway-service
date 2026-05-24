package resolver

import "net/http"

type ContextResolver struct{}

func NewContextResolver() *ContextResolver {
	return &ContextResolver{}
}

func (cr *ContextResolver) Resolve(r *http.Request, key string) (any, bool) {

	if v := r.Context().Value(key); v != nil {
		return v, true
	}

	return nil, false
}
