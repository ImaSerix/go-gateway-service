package resolver

import (
	"net/http"
)

type QueryResolver struct{}

func NewQueryResolver() *QueryResolver {
	return &QueryResolver{}
}

func (cr *QueryResolver) Resolve(r *http.Request, key string) (any, bool) {

	if v := r.URL.Query().Get(key); v != "" {
		return v, true
	}

	return nil, false
}
