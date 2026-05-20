package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
)

type ReverseBuilder struct {
	client *http.Client
}

func NewBuilder(client *http.Client) *ReverseBuilder {
	return &ReverseBuilder{
		client: client,
	}
}

func (rb *ReverseBuilder) Build(cfg config.Upstream) (pipeline.Proxy, error) {

	t, err := proxy.NewURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	m, err := proxy.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(string(t))
	if err != nil {
		return nil, err
	}

	p := httputil.NewSingleHostReverseProxy(url)

	p.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.Out.Method = string(m)
	}

	return p, nil
}
