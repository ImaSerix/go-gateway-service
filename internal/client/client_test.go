package client_test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/client"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}

	return u
}

func TestUpstream_DoRendersTargetAndMethod(t *testing.T) {
	var gotURL string
	var gotMethod string

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotURL = r.URL.String()
			gotMethod = r.Method
			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(http.NoBody),
			}, nil
		}),
	}

	target := mustParseURL(t, "http://policy.local/users/{query:id}")
	upstream := client.NewUpstreamClient(httpClient, renderer.NewRender(resolver.NewMultiResolver()), target, http.MethodPost)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/check?id=42", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := upstream.Do(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	if gotURL != "http://policy.local/users/42" {
		t.Fatalf("expected rendered url, got %q", gotURL)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected method %s, got %s", http.MethodPost, gotMethod)
	}
}
