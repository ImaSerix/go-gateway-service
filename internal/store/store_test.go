package store_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
	"github.com/ImaSerix/go-gateway-service/internal/store"
)

func TestStore_SaveFromResponse(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{"X-Token": []string{"header-token"}},
		Body:   io.NopCloser(strings.NewReader(`{"user_id":"42"}`)),
	}

	s := store.NewStore(map[string]string{
		"auth_token": "{header:X-Token}",
		"user_id":    "{body:user_id}",
	}, renderer.NewResponseRender(resolver.NewResponseMultiResolver()))

	ctx, err := s.Save(context.Background(), resp)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got := ctx.Value("auth_token"); got != "header-token" {
		t.Fatalf("expected auth token in context, got %v", got)
	}
	if got := ctx.Value("user_id"); got != "42" {
		t.Fatalf("expected user id in context, got %v", got)
	}
}

func TestStore_SaveNilResponse(t *testing.T) {
	s := store.NewStore(nil, renderer.NewResponseRender(resolver.NewResponseMultiResolver()))

	_, err := s.Save(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
