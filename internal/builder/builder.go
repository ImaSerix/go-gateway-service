package builder

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/transformer"
)

type EndpointBuilder struct {
	client *http.Client
}

func NewEndpointBuilder(c *http.Client) *EndpointBuilder {
	return &EndpointBuilder{
		client: c,
	}
}

func (b *EndpointBuilder) BuildEndpoint(cfg *config.RouteConfig) (*endpoint.Endpoint, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	path, err := endpoint.NewPath(cfg.Path)
	if err != nil {
		return nil, err
	}
	method, err := endpoint.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	if cfg.Upstream.Method == "" {
		cfg.Upstream.Method = cfg.Method
	}

	checks, err := checks.ChecksFactory(cfg.Checks, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	transformers, err := transformer.TransformersFactory(&cfg.Transform)
	if err != nil {
		return nil, err
	}

	target, err := proxy.NewURL(cfg.Upstream.URL)
	if err != nil {
		return nil, err
	}

	methodProxy, err := proxy.NewMethod(cfg.Upstream.Method)
	if err != nil {
		return nil, err
	}

	proxy := proxy.NewReverseProxy(target, methodProxy, b.client)

	e := endpoint.NewEndpoint(path, method, checks, transformers, proxy)
	return e, nil
}
