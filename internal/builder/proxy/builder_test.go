package proxy_test

import (
	"errors"
	"net/http"
	"testing"

	proxyBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
)

type rendererMock struct{}

func (rm *rendererMock) Render(s string, r *http.Request) (string, error) {
	return "", nil
}

func TestProxyBuilder(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.Upstream
		expErr error
	}{
		{
			name: "succes",
			cfg: config.Upstream{
				Host:   "nice.host",
				Scheme: "http",
				Path:   "",
				Method: "GET",
			},
			expErr: nil,
		},
		{
			name: "new url error",
			cfg: config.Upstream{
				Host:   "",
				Scheme: "http",
				Method: "GET",
			},
			expErr: proxy.ErrInvalidHost,
		},
		{
			name: "invalid method",
			cfg: config.Upstream{
				Host:   "nice.host",
				Scheme: "http",
				Path:   "",
				Method: "INVALID",
			},
			expErr: proxy.ErrInvalidMethod,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			b := proxyBuilder.NewBuilder(http.DefaultClient, &rendererMock{})

			p, err := b.Build(test.cfg)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, but got %v", test.expErr, err)
			}

			if err == nil && p == nil {
				t.Fatal("got no error, but proxy is nil")
			}
			if err != nil && p != nil {
				t.Fatal("got no error, but proxy is not nil")
			}

		})
	}
}
