package check_test

import (
	"context"
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
}

func (m *policyUpstreamMock) Do(r *http.Request) (*http.Response, error) {
	m.req = r
	return &http.Response{
		StatusCode: m.status,
		Body:       io.NopCloser(strings.NewReader(`{"token":"stored"}`)),
	}, nil
}

type policyStoreMock struct {
	called bool
}

func (m *policyStoreMock) Save(ctx context.Context, r *http.Response) (context.Context, error) {
	m.called = true
	return context.WithValue(ctx, "token", "stored"), nil
}

func TestPolicyCheck_ExecuteStoresResponseValues(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusAccepted}
	store := &policyStoreMock{}
	c := check.NewPolicyCheck([]pipeline.Transformer{policyTransformerMock{}}, upstream, store, http.StatusAccepted)

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
}

func TestPolicyCheck_ExecuteUnexpectedStatus(t *testing.T) {
	upstream := &policyUpstreamMock{status: http.StatusForbidden}
	store := &policyStoreMock{}
	c := check.NewPolicyCheck(nil, upstream, store, http.StatusAccepted)

	req, err := http.NewRequest(http.MethodGet, "http://gateway.local/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	if store.called {
		t.Fatal("expected store not to be called")
	}
}
