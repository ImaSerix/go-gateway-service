package proxy_test

import (
	"net/http"
	"testing"

	proxyBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
)

func TestProxyBuilder(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.Upstream
		expErr error
	}{
		{
			name: "succes",
			cfg: config.Upstream{
				URL:    "http://nice.url",
				Method: "GET",
			},
			expErr: nil,
		},
		{
			name: "new url error",
			cfg: config.Upstream{
				URL:    "http://",
				Method: "GET",
			},
			expErr: proxy.ErrEmptyHost,
		},
		{
			name: "invalid method",
			cfg: config.Upstream{
				URL:    "http://nice.url",
				Method: "INVALID",
			},
			expErr: proxy.ErrInvalidMethod,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			b := proxyBuilder.NewBuilder(http.DefaultClient)

			p, err := b.Build(test.cfg)
			if err != test.expErr {
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
