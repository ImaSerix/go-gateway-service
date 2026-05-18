package endpoint

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type CheckBuilder interface {
	Build(config.Check) (pipeline.Checker, error)
	BuildMany([]config.Check) ([]pipeline.Checker, error)
}

type MiddlewareBuilder interface {
	Build(config.Middleware) (pipeline.Middleware, error)
	BuildMany([]config.Middleware) ([]pipeline.Middleware, error)
}

type TransformBuilder interface {
	BuildMany(config.Transform) ([]pipeline.Transformer, error)
}

type ProxyBuilder interface {
	Build(config.Upstream) (pipeline.Proxy, error)
}
