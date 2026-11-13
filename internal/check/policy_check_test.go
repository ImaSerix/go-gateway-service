package check_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type policyTransformerMock struct {
	err error
}

func (m policyTransformerMock) Transform(r *http.Request) error {
	if m.err != nil {
		return m.err
	}
	r.Header.Set("X-Policy", "1")
	return nil
}

type policyUpstreamMock struct {
	req    *http.Request
	status int
	err    error
}

func (m *policyUpstreamMock) Do(r *http.Request) (*http.Response, error) {
	m.req = r
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.status,
		Body:       io.NopCloser(strings.NewReader(`{"token":"stored"}`)),
	}, nil
}

type policyStoreMock struct {
	called bool
	err    error
}

func (m *policyStoreMock) Save(ctx context.Context, r *http.Response) (context.Context, error) {
	m.called = true
	if m.err != nil {
		return ctx, m.err
	}
	return context.WithValue(ctx, "token", "stored"), nil
}

func TestPolicyCheck_ExecuteStoresResponseValues(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusAccepted}
	store := &policyStoreMock{}
	c := check.NewPolicyCheck([]pipeline.Transformer{policyTransformerMock{}}, upstream, store, http.StatusAccepted, true)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", strings.NewReader(`{"name":"Ada"}`))
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got := ctx.Value("token"); got != "stored" {
		t.Fatalf("expected stored token, got %v", got)
	}
	if !store.called {
		t.Fatal("expected store to be called")
	}
	if upstream.req.Header.Get("X-Policy") != "1" {
		t.Fatal("expected policy transform to run on upstream request")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("expected readable original body, got %v", err)
	}
	if string(body) != `{"name":"Ada"}` {
		t.Fatalf("expected original request body to be restored, got %q", string(body))
	}

	sent, err := io.ReadAll(upstream.req.Body)
	if err != nil {
		t.Fatalf("expected readable policy request body, got %v", err)
	}
	if string(sent) != `{"name":"Ada"}` {
		t.Fatalf("expected body to be forwarded to policy upstream, got %q", string(sent))
	}
}

func TestPolicyCheck_ExecuteUnexpectedStatus(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusForbidden}
	store := &policyStoreMock{}
	c := check.NewPolicyCheck(nil, upstream, store, http.StatusAccepted, false)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, check.ErrPolicyUnavailable) {
		t.Fatalf("expected %v, got %v", check.ErrPolicyUnavailable, err)
	}
	if store.called {
		t.Fatal("expected store not to be called")
	}
}

type policyBodyMock struct {
	read bool
}

func (m *policyBodyMock) Read(p []byte) (int, error) {
	m.read = true
	return 0, io.EOF
}

func (m *policyBodyMock) Close() error { return nil }

func TestPolicyCheck_ExecuteDoesNotForwardBody(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusOK}
	store := &policyStoreMock{}
	c := check.NewPolicyCheck(nil, upstream, store, http.StatusOK, false)

	body := &policyBodyMock{}
	req, err := http.NewRequest(http.MethodPost, "http://gateway.local/resource", body)
	if err != nil {
		t.Fatal(err)
	}
	req.ContentLength = 14
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")

	if _, err = c.Execute(context.Background(), req); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if body.read {
		t.Fatal("expected client body not to be read")
	}
	if req.Body != body || req.ContentLength != 14 {
		t.Fatal("expected client request to stay untouched")
	}
	if upstream.req.Body != http.NoBody || upstream.req.ContentLength != 0 {
		t.Fatal("expected policy request without body")
	}
	if upstream.req.Header.Get("Content-Type") != "" {
		t.Fatal("expected policy request without Content-Type")
	}
	if upstream.req.Header.Get("Authorization") != "Bearer token" {
		t.Fatal("expected client headers in policy request")
	}
}

func TestPolicyCheck_ExecuteErrors(t *testing.T) {
	tests := []struct {
		name     string
		upstream *policyUpstreamMock
		expErr   error
	}{
		{
			name:     "policy service answered 401",
			upstream: &policyUpstreamMock{status: http.StatusUnauthorized},
			expErr:   check.ErrUnauthorized,
		},
		{
			name:     "policy service answered 500",
			upstream: &policyUpstreamMock{status: http.StatusInternalServerError},
			expErr:   check.ErrPolicyUnavailable,
		},
		{
			name:     "policy service answered other 2xx",
			upstream: &policyUpstreamMock{status: http.StatusNoContent},
			expErr:   check.ErrPolicyUnavailable,
		},
		{
			name:     "policy service unreachable",
			upstream: &policyUpstreamMock{err: errors.New("dial error")},
			expErr:   check.ErrPolicyUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &policyStoreMock{}
			c := check.NewPolicyCheck(nil, test.upstream, store, http.StatusOK, false)

			req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = c.Execute(context.Background(), req)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected %v, got %v", test.expErr, err)
			}
			if store.called {
				t.Fatal("expected store not to be called")
			}
		})
	}
}

func TestPolicyCheck_ExecuteStoreError(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusOK}
	store := &policyStoreMock{err: errors.New("failed resolve")}
	c := check.NewPolicyCheck(nil, upstream, store, http.StatusOK, false)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, check.ErrPolicyUnavailable) {
		t.Fatalf("expected %v, got %v", check.ErrPolicyUnavailable, err)
	}
}

func TestPolicyCheck_ExecuteTransformError(t *testing.T) {
	transformErr := errors.New("failed transform")
	upstream := &policyUpstreamMock{status: http.StatusOK}
	c := check.NewPolicyCheck([]pipeline.Transformer{policyTransformerMock{err: transformErr}}, upstream, &policyStoreMock{}, http.StatusOK, false)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, transformErr) || errors.Is(err, check.ErrPolicyUnavailable) {
		t.Fatalf("expected plain transform error, got %v", err)
	}
	if upstream.req != nil {
		t.Fatal("expected policy upstream not to be called")
	}
}
