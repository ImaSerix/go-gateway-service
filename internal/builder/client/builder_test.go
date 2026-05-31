package client_test

import (
	"io"
	"net/http"
	"testing"

	clientbuilder "github.com/ImaSerix/go-gateway-service/internal/builder/client"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type rendererMock struct{}

func (m rendererMock) Render(s string, r *http.Request) (string, error) {
	return s, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestBuilder_Build(t *testing.T) {
	var gotURL string

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotURL = r.URL.String()
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(http.NoBody),
			}, nil
		}),
	}

	b := clientbuilder.NewBuilder(httpClient, rendererMock{})
	upstream, err := b.Build(config.Upstream{
		Scheme: "http",
		Host:   "policy.local",
		Path:   "/auth",
		Method: http.MethodPost,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := upstream.Do(req); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if gotURL != "http://policy.local/auth" {
		t.Fatalf("expected built upstream url, got %q", gotURL)
	}
}
