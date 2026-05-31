package check

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type QueryRequired struct {
	requiredQuery []string
}

func NewQueryRequired(cfg config.QueryRequiredCheck) (*QueryRequired, error) {

	if len(cfg.QueryParams) == 0 {
		return nil, ErrEmptyQuery
	}

	return &QueryRequired{
		requiredQuery: cfg.QueryParams,
	}, nil
}

func (c *QueryRequired) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	for _, q := range c.requiredQuery {
		if !r.URL.Query().Has(q) {
			return ctx, fmt.Errorf("%w: %s", ErrNoQueryParam, q)
		}
	}

	return ctx, nil
}
