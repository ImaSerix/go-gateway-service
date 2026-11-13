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

	m, err := client.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	// The caller must see the first response of the upstream: redirects are not followed.
	c := *b.client
	c.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return client.NewUpstreamClient(&c, b.render, u, string(m)), nil
}
