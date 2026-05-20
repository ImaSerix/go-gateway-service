package resolver

import (
	"net/http"
)

type QueryResolver struct{}

func NewQueryResolver() *QueryResolver {
	return &QueryResolver{}
}

func (cr *QueryResolver) Resolve(r *http.Request, key string) (any, bool) {

	_, n, err := getSourceAndKey(key)
	if err != nil {
		return nil, false
	}

	if v := r.URL.Query().Get(n); v != "" {
		return v, true
	}

	return nil, false
}
