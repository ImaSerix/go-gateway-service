package resolver

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type RouterResolver struct{}

func NewRouterResolver() *RouterResolver {
	return &RouterResolver{}
}

func (cr *RouterResolver) Resolve(r *http.Request, key string) (any, bool) {

	if v := chi.URLParam(r, key); v != "" {
		// chi routes by the escaped path when the request has one,
		// the parameter is then escaped too.
		if r.URL.RawPath != "" {
			if unescaped, err := url.PathUnescape(v); err == nil {
				v = unescaped
			}
		}
		return v, true
	}

	return nil, false
}
