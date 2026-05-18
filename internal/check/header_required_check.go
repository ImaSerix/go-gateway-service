package check

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const HeaderRequiredCheckType = "required_header"

type HeaderRequired struct {
	requiredHeaders []string
}

func NewHeaderRequiredCheck(cfg config.HeaderRequiredCheck) (*HeaderRequired, error) {

	if len(cfg.Header) == 0 {
		return nil, ErrEmptyHeaders
	}

	return &HeaderRequired{
		requiredHeaders: cfg.Header,
	}, nil
}

func (c *HeaderRequired) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	for _, h := range c.requiredHeaders {
		if v := r.Header.Get(h); v == "" {
			return ctx, fmt.Errorf("%w: %s", ErrNoHeader, h)
		}
	}

	return ctx, nil
}
