package builder_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	clientbuilder "github.com/ImaSerix/go-gateway-service/internal/builder/client"
	"github.com/ImaSerix/go-gateway-service/internal/builder/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/builder/handler"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/builder/render"
	resolverbuilder "github.com/ImaSerix/go-gateway-service/internal/builder/resolver"
	storebuilder "github.com/ImaSerix/go-gateway-service/internal/builder/store"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	resolverpkg "github.com/ImaSerix/go-gateway-service/internal/resolver"
	"gopkg.in/yaml.v3"
)

type stackRoundTripper struct {
	upstreamAuth string
}

func (rt *stackRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	switch r.URL.Host {
	case "policy.local":
		if r.Header.Get("X-Input") != "request-value" {
			return response(http.StatusForbidden, `{}`), nil
		}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     http.Header{"X-Token": []string{"stored-token"}},
			Body:       io.NopCloser(http.NoBody),
		}, nil
	case "api.local":
		rt.upstreamAuth = r.Header.Get("Authorization")
		return response(http.StatusOK, `{"ok":true}`), nil
	default:
		return response(http.StatusNotFound, `{}`), nil
	}
}

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestStack_PolicyStoreTransformProxy(t *testing.T) {
	transport := &stackRoundTripper{}
	httpClient := &http.Client{Transport: transport}

	checkRegistry := check.NewCheckRegistry()
	middlewareRegistry := middleware.NewMiddlewareRegistry()
	transformerRegistry := transformer.NewTransformerRegistry()

	requestRender := render.NewBuilder(resolverbuilder.NewMultiResolverBuilder().Build()).Build()
	responseRender := renderer.NewResponseRender(resolverpkg.NewResponseMultiResolver())

	transformerBuilder := transformer.NewBuilder(transformerRegistry)
	clientBuilder := clientbuilder.NewBuilder(httpClient, requestRender)
	storeBuilder := storebuilder.NewBuilder(responseRender)

	builder.RegisterChecks(checkRegistry, transformerBuilder, clientBuilder, storeBuilder)
	builder.RegisterTransformers(transformerRegistry, requestRender)
	builder.RegisterMiddlewares(middlewareRegistry)

	hb := handler.NewBuilder(
		middleware.NewBuilder(middlewareRegistry),
		endpoint.NewBuilder(
			check.NewBuilder(checkRegistry),
			transformerBuilder,
			middleware.NewBuilder(middlewareRegistry),
			proxy.NewBuilder(httpClient, requestRender),
		),
	)

	rawPolicyConfig := mustYAML(t, map[string]any{
		"transform": map[string]any{
			"header": map[string]any{
				"X-Input": "{header:X-Input}",
			},
		},
		"upstream": map[string]any{
			"scheme": "http",
			"host":   "policy.local",
			"path":   "/auth/{route:id}",
			"method": http.MethodGet,
		},
		"expected_status": http.StatusAccepted,
		"store": map[string]any{
			"token": "{header:X-Token}",
		},
	})

	h, err := hb.Build(config.Root{
		Routes: []config.Route{
			{
				Path:   "/resource/{id}",
				Method: http.MethodGet,
				Checks: []config.Check{
					{Type: "policy", Config: rawPolicyConfig},
				},
				Transforms: config.Transform{
					"header": mustYAML(t, map[string]any{
						"Authorization": "Bearer {context:token}",
					}),
				},
				Upstream: config.Upstream{
					Scheme: "http",
					Host:   "api.local",
					Path:   "/resource/{route:id}",
					Method: http.MethodGet,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/resource/42", nil)
	req.Header.Set("X-Input", "request-value")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if transport.upstreamAuth != "Bearer stored-token" {
		t.Fatalf("expected stored token in upstream auth header, got %q", transport.upstreamAuth)
	}
}

func mustYAML(t *testing.T, v any) yaml.Node {
	t.Helper()

	var node yaml.Node
	if err := node.Encode(v); err != nil {
		t.Fatal(err)
	}

	return node
}
