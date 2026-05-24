package client

import (
	"net/http"
	"net/url"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Upstream struct {
	client *http.Client
	render renderer.Renderer

	target *url.URL
	method string
}

func NewUpstreamClient(client *http.Client, render renderer.Renderer, target *url.URL, method string) *Upstream {
	return &Upstream{
		client: client,
		render: render,

		target: target,
		method: method,
	}
}

func (c *Upstream) Do(base *http.Request) (*http.Response, error) {

	r := base.Clone(base.Context())

	u, err := c.render.Render(c.target.String(), r)
	if err != nil {
		return nil, err
	}
	t, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	r.URL = t

	r.Method = c.method

	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
