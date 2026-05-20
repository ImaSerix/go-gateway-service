package resolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetSourceAndKey(t *testing.T) {
	_, _, err := getSourceAndKey("bad")
	if err != ErrInvalidKeyFormat { t.Fatalf("expected ErrInvalidKeyFormat, got %v", err) }
	_, _, err = getSourceAndKey(".id")
	if err != ErrEmptyKeySource { t.Fatalf("expected ErrEmptyKeySource, got %v", err) }
	_, _, err = getSourceAndKey("ctx.")
	if err != ErrEmptyKeyName { t.Fatalf("expected ErrEmptyKeyName, got %v", err) }
	s, n, err := getSourceAndKey(" query.id ")
	if err != nil || s != "query" || n != "id" { t.Fatalf("unexpected parse result: %s %s %v", s, n, err) }
}

func TestConcreteResolvers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1?id=7", nil)
	req.Header.Set("X-Req-ID", "abc")
	req = req.WithContext(context.WithValue(req.Context(), "token", "ctx-token"))

	ctxVal, ok := NewContextResolver().Resolve(req, "context.token")
	if !ok || ctxVal != "ctx-token" { t.Fatal("context resolver failed") }

	hVal, ok := NewHeaderResolver().Resolve(req, "header.X-Req-ID")
	if !ok || hVal != "abc" { t.Fatal("header resolver failed") }

	qVal, ok := NewQueryResolver().Resolve(req, "query.id")
	if !ok || qVal != "7" { t.Fatal("query resolver failed") }

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rVal, ok := NewRouterResolver().Resolve(req, "router.userID")
	if !ok || rVal != "42" { t.Fatal("router resolver failed") }
}

func TestMultiResolver(t *testing.T) {
	m := NewMultiResolver()
	m.Register("header", NewHeaderResolver())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Test", "ok")

	v, ok := m.Resolve(req, "header.X-Test")
	if !ok || v != "ok" { t.Fatal("expected resolved value") }

	_, ok = m.Resolve(req, "query.id")
	if ok { t.Fatal("expected false for unregistered resolver") }

	_, ok = m.Resolve(req, "bad")
	if ok { t.Fatal("expected false for invalid key") }
}
