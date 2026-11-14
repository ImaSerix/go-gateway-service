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

// A route parameter is the decoded value whether or not the request has an escaped path.
func TestRouterResolver_EscapedPath(t *testing.T) {
	tests := []struct {
		target string
		exp    string
	}{
		{target: "/users/42", exp: "42"},
		{target: "/users/42%23", exp: "42#"},
		{target: "/users/a%2Fb", exp: "a/b"},
		{target: "/users/a%252Fb", exp: "a%2Fb"},
		{target: "/users/%2e%2e", exp: ".."},
	}

	for _, test := range tests {
		t.Run(test.target, func(t *testing.T) {
			var got any

			router := chi.NewRouter()
			router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
				got, _ = NewRouterResolver().Resolve(r, "id")
			})
			router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, test.target, nil))

			if got != test.exp {
				t.Fatalf("expected %q, got %v", test.exp, got)
			}
		})
	}
}
