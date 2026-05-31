package handler

import (
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/go-chi/chi/v5"
)

type Builder struct {
	middlewareBuilder MiddlewareBuilder
	endpointBuilder   RouteBuilder
}

func NewBuilder(m MiddlewareBuilder, r RouteBuilder) *Builder {
	return &Builder{
		middlewareBuilder: m,
		endpointBuilder:   r,
	}
}

func (b *Builder) Build(cfg config.Root) (http.Handler, error) {

	middlewares, err := b.middlewareBuilder.BuildMany(cfg.Server.Middlewares)
	if err != nil {
		return nil, fmt.Errorf("global middleware build failed: %w", err)
	}

	endpoints, err := b.endpointBuilder.BuildMany(cfg.Routes)
	if err != nil {
		return nil, fmt.Errorf("route build failed: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middlewares...)

	for _, e := range endpoints {
		r.Method(e.Method(), e.Path(), e.Handler())
	}

	return r, nil
}
