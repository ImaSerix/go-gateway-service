package endpoint

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type Endpoint struct {
	path        Path
	method      Method
	check       []pipeline.Checker
	transformer []pipeline.Transformer
	middleware  []pipeline.Middleware
	proxy       pipeline.Proxy
}

func NewEndpoint(path Path, method Method, c []pipeline.Checker, t []pipeline.Transformer, p pipeline.Proxy, m []pipeline.Middleware) *Endpoint {
	return &Endpoint{
		path:        path,
		method:      method,
		check:       c,
		transformer: t,
		middleware:  m,
		proxy:       p,
	}
}

func (e *Endpoint) matchMethod(method string) bool {
	return e.method == Method(method)
}

func (e *Endpoint) Method() string {
	return string(e.method)
}

func (e *Endpoint) Path() string {
	return string(e.path)
}

func (e *Endpoint) Handler() http.Handler {
	return e
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !e.matchMethod(r.Method) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var h http.Handler = http.HandlerFunc(e.serve)

	for i := len(e.middleware) - 1; i >= 0; i-- {
		h = e.middleware[i](h)
	}

	h.ServeHTTP(w, r)
}

func (e *Endpoint) serve(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	for _, c := range e.check {

		var err error

		ctx, err = c.Execute(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}

	r = r.WithContext(ctx)

	for _, t := range e.transformer {
		err := t.Transform(r)
		if err != nil {
			http.Error(w, "transform failed", http.StatusInternalServerError)
			return
		}
	}

	e.proxy.ServeHTTP(w, r)
}
