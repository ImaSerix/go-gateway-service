package endpoint

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type ManyBuilder interface {
	BuildMany(raw yaml.Node) (pipeline.Middleware, error)
}

// Надо это переименовать, засунуть в другую папку, а может и осавить тут, но добавить билдеры как параметры
type EndpointBuilder struct {
	check      CheckBuilder
	transform  TransformBuilder
	middleware MiddlewareBuilder
	proxy      ProxyBuilder
}

func NewBuilder(ch CheckBuilder, t TransformBuilder, m MiddlewareBuilder, p ProxyBuilder) *EndpointBuilder {
	return &EndpointBuilder{
		check:      ch,
		transform:  t,
		middleware: m,
		proxy:      p,
	}
}

func (b *EndpointBuilder) Build(cfg config.Route) (pipeline.Endpoint, error) {

	check, err := b.check.BuildMany(cfg.Checks)
	if err != nil {
		return nil, err
	}

	transformer, err := b.transform.BuildMany(cfg.Transform)
	if err != nil {
		return nil, err
	}

	middleware, err := b.middleware.BuildMany(cfg.Middleware)
	if err != nil {
		return nil, err
	}

	if cfg.Upstream.Method == "" {
		cfg.Upstream.Method = cfg.Method
	}
	proxy, err := b.proxy.Build(cfg.Upstream)
	if err != nil {
		return nil, err
	}

	path, err := endpoint.NewPath(cfg.Path)
	if err != nil {
		return nil, err
	}
	method, err := endpoint.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}
	e := endpoint.NewEndpoint(path, method, check, transformer, proxy, middleware)

	return e, nil
}

func (b *EndpointBuilder) BuildMany(cfg []config.Route) ([]pipeline.Endpoint, error) {

	result := []pipeline.Endpoint{}

	for _, c := range cfg {
		e, err := b.Build(c)
		if err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, nil
}
