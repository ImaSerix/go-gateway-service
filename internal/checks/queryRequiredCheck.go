package checks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const QueryRequiredCheckType = "required_query "

type QueryRequiredCheck struct {
	requiredQueries []string
}

func NewQueryRequiredCheck(cfg *config.QueryRequiredCheckConfig) (*QueryRequiredCheck, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if len(cfg.Queries) == 0 {
		return nil, ErrEmptyQueries
	}

	return &QueryRequiredCheck{
		requiredQueries: cfg.Queries,
	}, nil
}

func (c *QueryRequiredCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	for _, q := range c.requiredQueries {
		if !r.URL.Query().Has(q) {
			return ctx, fmt.Errorf("%w: %s", ErrNoQueryParam, q)
		}
	}

	return ctx, nil
}
