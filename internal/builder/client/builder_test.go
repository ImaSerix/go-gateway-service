package client_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	clientbuilder "github.com/ImaSerix/go-gateway-service/internal/builder/client"
	"github.com/ImaSerix/go-gateway-service/internal/client"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type rendererMock struct{}

func (m rendererMock) Render(s string, r *http.Request) (string, error) {
	return s, nil
}

func (m rendererMock) RenderURL(s string, r *http.Request) (string, error) {
	return m.Render(s, r)
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

// Methods of an internal request are narrower than route methods: PATCH is not one of them.
func TestBuilder_Build_Methods(t *testing.T) {
	tests := []struct {
		method string
		expErr error
	}{
		{method: http.MethodGet},
		{method: http.MethodPost},
		{method: http.MethodPatch, expErr: client.ErrInvalidMethod},
		{method: http.MethodDelete, expErr: client.ErrInvalidMethod},
		{method: "", expErr: client.ErrInvalidMethod},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			b := clientbuilder.NewBuilder(&http.Client{}, rendererMock{})

			_, err := b.Build(config.Upstream{Scheme: "http", Host: "policy.local", Method: test.method})
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, got %v", test.expErr, err)
			}
		})
	}
}

func TestBuilder_Build_DoesNotFollowRedirects(t *testing.T) {
	calls := 0

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			return &http.Response{
				StatusCode: http.StatusFound,
				Header:     http.Header{"Location": []string{"/redirected"}},
				Body:       io.NopCloser(http.NoBody),
			}, nil
		}),
	}

	upstream, err := clientbuilder.NewBuilder(httpClient, rendererMock{}).Build(config.Upstream{
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

	resp, err := upstream.Do(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != http.StatusFound || calls != 1 {
		t.Fatalf("expected the first response (302) after one call, got %d after %d calls", resp.StatusCode, calls)
	}
	if httpClient.CheckRedirect != nil {
		t.Fatal("expected the shared http client to stay unchanged")
	}
}
