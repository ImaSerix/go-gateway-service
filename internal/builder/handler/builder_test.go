package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder/handler"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type middlewareBuilderMock struct {
	called bool
	result []pipeline.Middleware
	err    error
}

func (m *middlewareBuilderMock) BuildMany(_ []config.Middleware) ([]pipeline.Middleware, error) {
	m.called = true
	return m.result, m.err
}

type endpointBuilderMock struct {
	called bool
	result []pipeline.Endpoint
	err    error
}

func (e *endpointBuilderMock) BuildMany(_ []config.Route) ([]pipeline.Endpoint, error) {
	e.called = true
	return e.result, e.err
}

type endpointMock struct {
	path    string
	method  string
	handler http.Handler
}

func (e endpointMock) Path() string           { return e.path }
func (e endpointMock) Method() string         { return e.method }
func (e endpointMock) Handler() http.Handler  { return e.handler }

func TestHandlerBuilder_Build(t *testing.T) {
	endpointCalled := false
	endpoint := endpointMock{path: "/v1", method: "GET", handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		endpointCalled = true
		w.WriteHeader(http.StatusNoContent)
	})}

	mb := &middlewareBuilderMock{result: []pipeline.Middleware{func(next http.Handler) http.Handler { return next }}}
	eb := &endpointBuilderMock{result: []pipeline.Endpoint{endpoint}}
	b := handler.NewBuilder(mb, eb)

	h, err := b.Build(config.Root{})
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, but got %d", http.StatusNoContent, rr.Code)
	}
	if !endpointCalled {
		t.Fatal("expected endpoint handler to be called")
	}
	if !mb.called || !eb.called {
		t.Fatal("expected both builders to be called")
	}
}

func TestHandlerBuilder_Build_Errors(t *testing.T) {
	badErr := errors.New("bad")

	t.Run("middleware builder error", func(t *testing.T) {
		mb := &middlewareBuilderMock{err: badErr}
		eb := &endpointBuilderMock{}
		b := handler.NewBuilder(mb, eb)

		h, err := b.Build(config.Root{})
		if !errors.Is(err, badErr) {
			t.Fatalf("expected %v, got %v", badErr, err)
		}
		if h != nil {
			t.Fatal("expected nil handler")
		}
		if eb.called {
			t.Fatal("endpoint builder should not be called")
		}
	})

	t.Run("endpoint builder error", func(t *testing.T) {
		mb := &middlewareBuilderMock{}
		eb := &endpointBuilderMock{err: badErr}
		b := handler.NewBuilder(mb, eb)

		h, err := b.Build(config.Root{})
		if !errors.Is(err, badErr) {
			t.Fatalf("expected %v, got %v", badErr, err)
		}
		if h != nil {
			t.Fatal("expected nil handler")
		}
	})
}
