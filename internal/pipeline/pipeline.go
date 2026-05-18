package pipeline

import (
	"context"
	"net/http"
)

type Checker interface {
	Execute(ctx context.Context, r *http.Request) (context.Context, error)
}

type Transformer interface {
	Transform(ctx context.Context, r *http.Request) error
}

type Proxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Middleware = func(http.Handler) http.Handler

type Endpoint interface {
	Method() string
	Path() string
	Handler() http.Handler
}
