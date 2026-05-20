package resolver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouterResolver struct{}

func NewRouterResolver() *RouterResolver {
	return &RouterResolver{}
}

func (cr *RouterResolver) Resolve(r *http.Request, key string) (any, bool) {

	_, n, err := getSourceAndKey(key)
	if err != nil {
		return nil, false
	}

	if v := chi.URLParam(r, n); v != "" {
		return v, true
	}

	return nil, false
}
