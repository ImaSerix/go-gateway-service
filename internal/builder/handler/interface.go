package handler

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type MiddlewareBuilder interface {
	BuildMany([]config.Middleware) ([]pipeline.Middleware, error)
}

type RouteBuilder interface {
	BuildMany([]config.Route) ([]pipeline.Endpoint, error)
}
