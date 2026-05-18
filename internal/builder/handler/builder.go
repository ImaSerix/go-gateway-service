package handler

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/go-chi/chi/v5"
)

type Builder struct {
	middlewareBuilder MiddlewareBuilder
	endpointBuilder   EndpointBuilder
}

func NewBuilder(m MiddlewareBuilder, e EndpointBuilder) *Builder {
	return &Builder{
		middlewareBuilder: m,
		endpointBuilder:   e,
	}
}

func (b *Builder) Build(cfg config.Root) (http.Handler, error) {

	middlewares, err := b.middlewareBuilder.BuildMany(cfg.Server.Middleware)
	if err != nil {
		return nil, err
	}

	endpoints, err := b.endpointBuilder.BuildMany(cfg.Routes)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Use(middlewares...)

	for _, e := range endpoints {
		r.Method(e.Method(), e.Path(), e.Handler())
	}

	return r, nil
}
