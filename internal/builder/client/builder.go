package client

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/client"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Builder struct {
	client *http.Client
	render renderer.Renderer
}

func NewBuilder(c *http.Client, render renderer.Renderer) *Builder {
	return &Builder{
		client: c,
		render: render,
	}
}

func (b *Builder) Build(cfg config.Upstream) (*client.Upstream, error) {

	u, err := proxy.MakeURL(cfg.Scheme, cfg.Host, cfg.Path)
	if err != nil {
		return nil, err
	}

	m, err := proxy.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	return client.NewUpstreamClient(b.client, b.render, u, string(m)), nil
}
