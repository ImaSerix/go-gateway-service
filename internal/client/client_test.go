package client_test

import (
	"errors"
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

// A value from the request cannot change the path or the query of the internal request.
func TestUpstream_DoEscapesValues(t *testing.T) {
	tests := []struct {
		name   string
		target string
		query  string

		expURL string
		expErr error
	}{
		{
			name:   "path value",
			target: "http://policy.local/users/{query:id}/check",
			query:  "id=42%3Fadmin%3Dtrue%23",
			expURL: "http://policy.local/users/42%3Fadmin=true%23/check",
		},
		{
			name:   "slashes in path value",
			target: "http://policy.local/users/{query:id}",
			query:  "id=..%2Fadmin%2Fdelete",
			expURL: "http://policy.local/users/..%2Fadmin%2Fdelete",
		},
		{
			name:   "query value",
			target: "http://policy.local/check?locale={query:locale}",
			query:  "locale=lv%26admin%3Dtrue",
			expURL: "http://policy.local/check?locale=lv%26admin%3Dtrue",
		},
		{
			name:   "dot segment value",
			target: "http://policy.local/users/{query:id}/check",
			query:  "id=..",
			expErr: renderer.ErrDotSegment,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotURL string
			calls := 0

			httpClient := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					calls++
					gotURL = r.URL.String()
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(http.NoBody),
					}, nil
				}),
			}

			upstream := client.NewUpstreamClient(httpClient, renderer.NewRender(resolver.NewMultiResolver()), mustParseURL(t, test.target), http.MethodPost)

			req, err := http.NewRequest(http.MethodGet, "http://gateway.local/check?"+test.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = upstream.Do(req)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, got %v", test.expErr, err)
			}
			if test.expErr != nil {
				if calls != 0 {
					t.Fatalf("expected no request, got %d", calls)
				}
				return
			}
			if gotURL != test.expURL {
				t.Fatalf("expected url %q, got %q", test.expURL, gotURL)
			}
		})
	}
}
