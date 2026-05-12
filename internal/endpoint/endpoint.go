package endpoint

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type Endpoint struct {
	path         Path
	method       Method
	checks       []pipeline.Checker
	transformers []pipeline.Transformer
	proxy        pipeline.Proxy
}

func NewEndpoint(path Path, method Method, checks []pipeline.Checker, transformers []pipeline.Transformer, proxy pipeline.Proxy) *Endpoint {
	return &Endpoint{
		path:         path,
		method:       method,
		checks:       checks,
		transformers: transformers,
		proxy:        proxy,
	}
}

func (e *Endpoint) matchMethod(method string) bool {
	return e.method == Method(method)
}

func (e *Endpoint) Pattern() string {
	return string(e.method) + " " + string(e.path)
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !e.matchMethod(r.Method) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	for _, c := range e.checks {

		var err error

		ctx, err = c.Execute(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}

	r = r.WithContext(ctx)

	for _, t := range e.transformers {
		err := t.Transform(r.Context(), r)
		if err != nil {
			http.Error(w, "transform failed", http.StatusInternalServerError)
			return
		}
	}

	e.proxy.ServeHTTP(w, r)
}
