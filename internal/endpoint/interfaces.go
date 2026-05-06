package endpoint

import (
	"context"
	"net/http"
)

type Check interface {
	Execute(ctx context.Context, r *http.Request) (context.Context, error)
}
