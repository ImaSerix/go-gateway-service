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

	if v := chi.URLParam(r, key); v != "" {
		return v, true
	}

	return nil, false
}
