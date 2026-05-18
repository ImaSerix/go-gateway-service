package check

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type Inject struct {
	ctx map[string]any
}

func NewInject(cfg config.InjectCheck) (*Inject, error) {

	if len(cfg.Ctx) == 0 {
		return nil, ErrEmptyInjectContext
	}

	return &Inject{
		ctx: cfg.Ctx,
	}, nil
}

func (c *Inject) Execute(ctx context.Context, r *http.Request) (context.Context, error) {

	if r == nil {
		return ctx, ErrNilRequest
	}

	for key, value := range c.ctx {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx, nil
}
