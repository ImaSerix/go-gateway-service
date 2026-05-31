package store_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	storebuilder "github.com/ImaSerix/go-gateway-service/internal/builder/store"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

func TestBuilder_Build(t *testing.T) {
	b := storebuilder.NewBuilder(renderer.NewResponseRender(resolver.NewResponseMultiResolver()))

	s, err := b.Build(config.Store{"token": "{header:X-Token}"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	resp := &http.Response{
		Header: http.Header{"X-Token": []string{"abc"}},
		Body:   io.NopCloser(strings.NewReader(`{}`)),
	}

	ctx, err := s.Save(t.Context(), resp)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got := ctx.Value("token"); got != "abc" {
		t.Fatalf("expected token in context, got %v", got)
	}
}
