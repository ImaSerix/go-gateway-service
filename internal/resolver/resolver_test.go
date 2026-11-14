package resolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestConcreteResolvers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1?id=7", nil)
	req.Header.Set("X-Req-ID", "abc")
	req = req.WithContext(context.WithValue(req.Context(), "token", "ctx-token"))

	ctxVal, ok := NewContextResolver().Resolve(req, "token")
	if !ok || ctxVal != "ctx-token" {
		t.Fatal("context resolver failed")
	}

	hVal, ok := NewHeaderResolver().Resolve(req, "X-Req-ID")
	if !ok || hVal != "abc" {
		t.Fatal("header resolver failed")
	}

	qVal, ok := NewQueryResolver().Resolve(req, "id")
	if !ok || qVal != "7" {
		t.Fatal("query resolver failed")
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rVal, ok := NewRouterResolver().Resolve(req, "userID")
	if !ok || rVal != "42" {
		t.Fatal("router resolver failed")
	}
}

func TestMultiResolver(t *testing.T) {
	m := NewMultiResolver()
	m.Register("header", NewHeaderResolver())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Test", "ok")

	v, ok := m.Resolve(req, "header", "X-Test")
	if !ok || v != "ok" {
		t.Fatal("expected resolved value")
	}

	_, ok = m.Resolve(req, "query", "id")
	if ok {
		t.Fatal("expected false for unregistered resolver")
	}

}
