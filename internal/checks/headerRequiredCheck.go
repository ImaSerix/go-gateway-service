package checks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const HeaderRequiredCheckType = "required_header"

type HeaderRequiredCheck struct {
	requiredHeaders []string
}

func NewHeaderRequiredCheck(cfg *config.HeaderRequiredCheckConfig) (*HeaderRequiredCheck, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if len(cfg.Headers) == 0 {
		return nil, ErrEmptyHeaders
	}

	return &HeaderRequiredCheck{
		requiredHeaders: cfg.Headers,
	}, nil
}

func (c *HeaderRequiredCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
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
