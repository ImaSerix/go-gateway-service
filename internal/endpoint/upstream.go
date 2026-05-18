package endpoint

import "github.com/ImaSerix/go-gateway-service/internal/config"

type Upstream struct {
	URL    URL
	Method Method
}

func NewUpstream(url URL, method Method) *Upstream {
	return &Upstream{
		URL:    url,
		Method: method,
	}
}

func NewUpstreamFromConfig(cfg *config.Upstream) (*Upstream, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	url, err := NewURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	method, err := NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	return &Upstream{
		URL:    url,
		Method: method,
	}, nil
}
