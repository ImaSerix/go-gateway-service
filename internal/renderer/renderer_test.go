package renderer_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

func TestRender_RequestHeaderWithHyphen(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Token", "abc")

	r := renderer.NewRender(resolver.NewMultiResolver())
	got, err := r.Render("Bearer {header:X-Token}", req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got != "Bearer abc" {
		t.Fatalf("expected rendered header, got %q", got)
	}
}

func TestResponseRender_HeaderAndTopLevelBody(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{"X-Token": []string{"abc"}},
		Body:   io.NopCloser(strings.NewReader(`{"token":"body-token","nested":{"ignored":true}}`)),
	}

	r := renderer.NewResponseRender(resolver.NewResponseMultiResolver())
	got, err := r.Render("{header:X-Token}:{body:token}", resp)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got != "abc:body-token" {
		t.Fatalf("expected rendered response values, got %q", got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected readable body, got %v", err)
	}
	if string(body) != `{"token":"body-token","nested":{"ignored":true}}` {
		t.Fatalf("expected response body to be restored, got %q", string(body))
	}
}
