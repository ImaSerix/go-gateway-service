package renderer_test

import (
	"context"
	"errors"
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

func TestRender_RenderURL(t *testing.T) {
	tests := []struct {
		name     string
		template string
		target   string

		exp    string
		expErr error
	}{
		{
			name:     "regular values are not changed",
			template: "http://api.local/users/{query:id}/items?name={query:name}",
			target:   "/?id=3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e&name=Abc-123_x",
			exp:      "http://api.local/users/3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e/items?name=Abc-123_x",
		},
		{
			name:     "path value is escaped as a path segment",
			template: "http://api.local/users/{query:id}/items",
			target:   "/?id=a%2Fb%3Fc%23d%25e%20f",
			exp:      "http://api.local/users/a%2Fb%3Fc%23d%25e%20f/items",
		},
		{
			name:     "query value is escaped as a query value",
			template: "http://api.local/search?name={query:name}&fixed=1",
			target:   "/?name=a%26b%3Dc%23d%25e%20f",
			exp:      "http://api.local/search?name=a%26b%3Dc%23d%25e+f&fixed=1",
		},
		{
			name:     "header value",
			template: "http://api.local/users/{header:X-Name}?name={header:X-Name}",
			target:   "/",
			exp:      "http://api.local/users/a%2Fb&c?name=a%2Fb%26c",
		},
		{
			name:     "template text is not escaped",
			template: "http://api.local/a%2Fb/../{query:id}?x=1&y=%26",
			target:   "/?id=42",
			exp:      "http://api.local/a%2Fb/../42?x=1&y=%26",
		},
		{
			name:     "one dot in path",
			template: "http://api.local/users/{query:id}/items",
			target:   "/?id=.",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "two dots in path",
			template: "http://api.local/users/{query:id}",
			target:   "/?id=..",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "dot segment made of two values",
			template: "http://api.local/users/{query:a}{query:b}/items",
			target:   "/?a=.&b=.",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "dot segment made of template text and a value",
			template: "http://api.local/users/.{query:a}/items",
			target:   "/?a=.",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "one dot made of template text and an empty value",
			template: "http://api.local/users/.{context:empty}/items",
			target:   "/",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "two dots made of template text and an empty value",
			template: "http://api.local/users/..{context:empty}",
			target:   "/",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "empty value next to other template text",
			template: "http://api.local/users/v.{context:empty}/items",
			target:   "/",
			exp:      "http://api.local/users/v./items",
		},
		{
			name:     "dot value before a missing value",
			template: "http://api.local/users/{query:id}/{query:missing}",
			target:   "/?id=..",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "dot value after a missing value",
			template: "http://api.local/users/{query:missing}/{query:id}",
			target:   "/?id=..",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "dot value with a missing value in query",
			template: "http://api.local/users/{query:id}?x={query:missing}",
			target:   "/?id=..",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "made dot segment with a missing value in another segment",
			template: "http://api.local/users/{query:missing}/.{context:empty}",
			target:   "/",
			expErr:   renderer.ErrDotSegment,
		},
		{
			name:     "missing value next to a dot in the same segment",
			template: "http://api.local/users/.{query:missing}",
			target:   "/",
			expErr:   renderer.ErrFailedResolve,
		},
		{
			name:     "dots in query",
			template: "http://api.local/search?dir={query:dir}",
			target:   "/?dir=..",
			exp:      "http://api.local/search?dir=..",
		},
		{
			name:     "missing value",
			template: "http://api.local/users/{query:id}",
			target:   "/",
			expErr:   renderer.ErrFailedResolve,
		},
		{
			name:     "missing value in query",
			template: "http://api.local/search?name={query:name}",
			target:   "/",
			expErr:   renderer.ErrFailedResolve,
		},
		{
			name:     "brackets are not closed",
			template: "http://api.local/users/{query:id",
			target:   "/?id=42",
			expErr:   renderer.ErrInvalidTemplate,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			req.Header.Set("X-Name", "a/b&c")
			req = req.WithContext(context.WithValue(req.Context(), "empty", ""))

			got, err := renderer.NewRender(resolver.NewMultiResolver()).RenderURL(test.template, req)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, got %v", test.expErr, err)
			}
			if got != test.exp {
				t.Fatalf("expected %q, got %q", test.exp, got)
			}
		})
	}
}

// Render itself does not escape: headers, query parameters and body fields
// of transforms get the value as it is.
func TestRender_DoesNotEscape(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?id=a%2Fb%3Fc&dir=..", nil)

	got, err := renderer.NewRender(resolver.NewMultiResolver()).Render("{query:id} {query:dir}", req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != "a/b?c .." {
		t.Fatalf("expected raw values, got %q", got)
	}
}
