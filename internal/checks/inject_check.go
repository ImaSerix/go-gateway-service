package checks

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const InjectCheckType = "inject"

type InjectCheck struct {
	ctx map[string]any
}

func NewInjectCheck(cfg *config.InjectCheckConfig) (*InjectCheck, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if len(cfg.Ctx) == 0 {
		return nil, ErrEmptyInjectContext
	}

	return &InjectCheck{
		ctx: cfg.Ctx,
	}, nil
}

func (c *InjectCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {

	if r == nil {
		return ctx, ErrNilRequest
	}

	for key, value := range c.ctx {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx, nil
}
