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
