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
				URL:    "https://jsonplaceholder.typicode.com/users",
				Method: "GET",
			},
			expError: nil,
		},
		{
			name:     "nil config",
			cfg:      nil,
			expError: endpoint.ErrInvalidConfig,
		},
		{
			name: "invalid url",
			cfg: &config.Upstream{
				URL:    "https:///users",
				Method: "GET",
			},
			expError: endpoint.ErrEmptyHost,
		},
		{
			name: "invalid method",
			cfg: &config.Upstream{
				URL:    "https://jsonplaceholder.typicode.com/users",
				Method: "INVALID",
			},
			expError: endpoint.ErrInvalidMethod,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			end, err := endpoint.NewUpstreamFromConfig(test.cfg)
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
				if end.URL != endpoint.URL(test.cfg.URL) {
					t.Fatalf("expected upstream url %s, but got %s", test.cfg.URL, end.URL)
				}
				if end.Method != endpoint.Method(test.cfg.Method) {
					t.Fatalf("expected upstream method %s, but got %s", test.cfg.Method, end.Method)
				}
			}
		})
	}
}
