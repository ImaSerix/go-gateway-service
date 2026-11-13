package endpoint

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/httperror"
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
		httperror.Write(w, http.StatusMethodNotAllowed, httperror.CodeMethodNotAllowed, "method not allowed")
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
			slog.Warn("check failed", "method", r.Method, "path", r.URL.Path, "err", err)
			writeCheckError(w, err)
			return
		}
	}

	r = r.WithContext(ctx)

	for _, t := range e.transformer {
		err := t.Transform(r)
		if err != nil {
			slog.Error("transform failed", "method", r.Method, "path", r.URL.Path, "err", err)
			httperror.Write(w, http.StatusInternalServerError, httperror.CodeInternalError, "internal error")
			return
		}
	}

	e.proxy.ServeHTTP(w, r)
}

// writeCheckError maps a failed check to the client response.
// The error text itself is never sent to the client.
func writeCheckError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, check.ErrUnauthorized):
		httperror.Write(w, http.StatusUnauthorized, httperror.CodeUnauthorized, "authentication required")
	case errors.Is(err, check.ErrPolicyUnavailable):
		httperror.Write(w, http.StatusServiceUnavailable, httperror.CodeAuthUnavailable, "authentication service unavailable")
	case errors.Is(err, check.ErrForbidden), errors.Is(err, check.ErrNoHeader), errors.Is(err, check.ErrNoQueryParam):
		httperror.Write(w, http.StatusForbidden, httperror.CodeForbidden, "access denied")
	default:
		httperror.Write(w, http.StatusInternalServerError, httperror.CodeInternalError, "internal error")
	}
}
