package endpoint_test

import (
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
)

func TestUpstreamFromConfig(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Upstream
		expError error
	}{
		{
			name: "success",
			cfg: &config.Upstream{
				
			},
			},
			expError: nil,
		},
		{
			name:     "nil config",
			cfg:      nil,
			expError: endpoint.ErrInvalidConfig,
		},
		{
			name: "invalid path",
			cfg: &config.RouteConfig{
				Path:   "",
				Method: "GET",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "GET",
				},
			},
			expError: endpoint.ErrEmptyPath,
		},
		{
			name: "invalid method",
			cfg: &config.RouteConfig{
				Path:   "/nice/path",
				Method: "INVALID",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "GET",
				},
			},
			expError: endpoint.ErrInvalidMethod,
		},
		{
			name: "invalid URL",
			cfg: &config.RouteConfig{
				Path:   "/nice/path",
				Method: "GET",
				Upstream: config.Upstream{
					URL:    "https:///users",
					Method: "GET",
				},
			},
			expError: endpoint.ErrEmptyHost,
		},
		{
			name: "invalid target method",
			cfg: &config.RouteConfig{
				Path:   "/nice/path",
				Method: "GET",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "INVALID",
				},
			},
			expError: endpoint.ErrInvalidMethod,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			end, err := endpoint.NewEndpointFromConfig(test.cfg)
			if err != test.expError {
				t.Fatalf("expected error %v, but got %v", test.expError, err)
			}
			if end != nil && err != nil {
				t.Fatal("got erorr, but endpoint is not nil")
			}
			if end == nil && err == nil {
				t.Fatal("got nil endpoint, but no error")
			}
			if end != nil && test.cfg != nil {
				if end.Path != endpoint.Path(test.cfg.Path) {
					t.Fatalf("expected path %s, but got %s", test.cfg.Path, end.Path)
				}
				if end.Method != endpoint.Method(test.cfg.Method) {
					t.Fatalf("expected method %s, but got %s", test.cfg.Method, end.Method)
				}
				if end.Upstream == nil {
					t.Fatal("expected non-nil upstream, but got nil")
				}
				if end.Upstream.URL != endpoint.URL(test.cfg.Upstream.URL) {
					t.Fatalf("expected upstream url %s, but got %s", test.cfg.Upstream.URL, end.Upstream.URL)
				}
				if end.Upstream.Method != endpoint.Method(test.cfg.Upstream.Method) {
					t.Fatalf("expected upstream method %s, but got %s", test.cfg.Upstream.Method, end.Upstream.Method)
				}
			}
		})
	}
}
