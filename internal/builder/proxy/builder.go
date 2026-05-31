package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type ReverseBuilder struct {
	client *http.Client
	render renderer.Renderer
}

func NewBuilder(client *http.Client, render renderer.Renderer) *ReverseBuilder {
	return &ReverseBuilder{
		client: client,
		render: render,
	}
}

func (rb *ReverseBuilder) Build(cfg config.Upstream) (pipeline.Proxy, error) {

	u, err := proxy.MakeURL(cfg.Scheme, cfg.Host, cfg.Path)
	if err != nil {
		return nil, err
	}

	m, err := proxy.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	var transport http.RoundTripper
	if rb.client != nil {
		transport = rb.client.Transport
	}

	p := &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {

			// TODO: думаю тут будет ошибка
			renderedURL, err := rb.render.Render(rawTarget(u), pr.In)
			if err != nil {
				return
			}

			u, err := url.Parse(renderedURL)
			if err != nil {
				return
			}

			pr.SetURL(u)
			pr.Out.Method = string(m)
		},
	}

	return p, nil
}

func rawTarget(u *url.URL) string {
	raw := u.Scheme + "://" + u.Host + u.Path
	if u.RawQuery != "" {
		raw += "?" + u.RawQuery
	}

	return raw
}
