package proxy_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

type failingTransport struct{}

func (failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("dial tcp 10.0.0.1:8080: connection refused")
}

func TestProxyBuilder_UpstreamUnavailable(t *testing.T) {
	b := proxyBuilder.NewBuilder(&http.Client{Transport: failingTransport{}}, &rendererMock{})

	p, err := b.Build(config.Upstream{Host: "nice.host", Scheme: "http", Method: "GET"})
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	w := httptest.NewRecorder()
	p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users", nil))

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, but got %d", http.StatusBadGateway, w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"upstream_unavailable"`) {
		t.Fatalf("expected upstream_unavailable code, but got %s", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "10.0.0.1") {
		t.Fatalf("internal error text leaked to client: %s", w.Body.String())
	}
}
