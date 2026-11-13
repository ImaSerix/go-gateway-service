package builder_test

import (
	"errors"
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

	hb := newStackHandlerBuilder(httpClient)

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

func newStackHandlerBuilder(httpClient *http.Client) *handler.Builder {
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

	return handler.NewBuilder(
		middleware.NewBuilder(middlewareRegistry),
		endpoint.NewBuilder(
			check.NewBuilder(checkRegistry),
			transformerBuilder,
			middleware.NewBuilder(middlewareRegistry),
			proxy.NewBuilder(httpClient, requestRender),
		),
	)
}

func mustYAML(t *testing.T, v any) yaml.Node {
	t.Helper()

	var node yaml.Node
	if err := node.Encode(v); err != nil {
		t.Fatal(err)
	}

	return node
}

type authStackRoundTripper struct {
	authErr error

	authCalls   int
	authMethod  string
	authPath    string
	authBody    string
	authHeader  string
	authUserID  string
	apiCalls    int
	apiPath     string
	apiQuery    string
	apiUserID   string
	apiBody     string
	apiAuthSeen string
}

func (rt *authStackRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := ""
	if r.Body != nil {
		raw, _ := io.ReadAll(r.Body)
		body = string(raw)
	}

	switch r.URL.Host {
	case "auth.local":
		rt.authCalls++
		rt.authMethod = r.Method
		rt.authPath = r.URL.Path
		rt.authBody = body
		if r.URL.Path == "/redirected" {
			return response(http.StatusOK, `{"id":"user-1"}`), nil
		}
		rt.authHeader = r.Header.Get("Authorization")
		rt.authUserID = strings.Join(r.Header.Values("X-User-ID"), ",")
		if rt.authErr != nil {
			return nil, rt.authErr
		}
		switch rt.authHeader {
		case "Bearer redirect":
			return &http.Response{
				StatusCode: http.StatusFound,
				Header:     http.Header{"Location": []string{"/redirected"}},
				Body:       io.NopCloser(http.NoBody),
			}, nil
		case "Bearer no-id":
			return response(http.StatusOK, `{"username":"alice"}`), nil
		case "Bearer good":
			return response(http.StatusOK, `{"id":"user-1","username":"alice","roles":["user"]}`), nil
		case "Bearer broken-auth":
			return response(http.StatusInternalServerError, `internal error`), nil
		default:
			return response(http.StatusUnauthorized, `invalid token`), nil
		}
	case "api.local":
		rt.apiCalls++
		rt.apiPath = r.URL.Path
		rt.apiQuery = r.URL.RawQuery
		rt.apiUserID = strings.Join(r.Header.Values("X-User-ID"), ",")
		rt.apiBody = body
		rt.apiAuthSeen = r.Header.Get("Authorization")
		return response(http.StatusCreated, `{"ok":true}`), nil
	default:
		return response(http.StatusNotFound, `{}`), nil
	}
}

// Token check through a policy request without the client body: the verified
// user id replaces X-User-ID from the client, path, query and body reach the service.
func TestStack_TokenCheckSetsUserID(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		target  string
		body    string
		token   string
		authErr error

		expCode    int
		expErrCode string
		expAPIPath string
		expQuery   string
	}{
		{
			name:       "GET without body",
			method:     http.MethodGet,
			target:     "/chats/7/messages?limit=5&before=10",
			token:      "Bearer good",
			expCode:    http.StatusCreated,
			expAPIPath: "/chats/7/messages",
			expQuery:   "limit=5&before=10",
		},
		{
			name:       "POST with body",
			method:     http.MethodPost,
			target:     "/chats/7/messages",
			body:       `{"text":"hello"}`,
			token:      "Bearer good",
			expCode:    http.StatusCreated,
			expAPIPath: "/chats/7/messages",
		},
		{
			name:       "PATCH with body",
			method:     http.MethodPatch,
			target:     "/chats/7/messages",
			body:       `{"text":"edited"}`,
			token:      "Bearer good",
			expCode:    http.StatusCreated,
			expAPIPath: "/chats/7/messages",
		},
		{
			name:       "no token",
			method:     http.MethodPost,
			target:     "/chats/7/messages",
			body:       `{"text":"hello"}`,
			expCode:    http.StatusUnauthorized,
			expErrCode: "unauthorized",
		},
		{
			name:       "invalid token",
			method:     http.MethodGet,
			target:     "/chats/7/messages",
			token:      "Bearer bad",
			expCode:    http.StatusUnauthorized,
			expErrCode: "unauthorized",
		},
		{
			name:       "auth answered with unexpected status",
			method:     http.MethodGet,
			target:     "/chats/7/messages",
			token:      "Bearer broken-auth",
			expCode:    http.StatusServiceUnavailable,
			expErrCode: "auth_unavailable",
		},
		{
			name:       "auth answered with redirect",
			method:     http.MethodGet,
			target:     "/chats/7/messages",
			token:      "Bearer redirect",
			expCode:    http.StatusServiceUnavailable,
			expErrCode: "auth_unavailable",
		},
		{
			name:       "auth answered 200 without id",
			method:     http.MethodGet,
			target:     "/chats/7/messages",
			token:      "Bearer no-id",
			expCode:    http.StatusServiceUnavailable,
			expErrCode: "auth_unavailable",
		},
		{
			name:       "auth unreachable",
			method:     http.MethodGet,
			target:     "/chats/7/messages",
			token:      "Bearer good",
			authErr:    errors.New("dial tcp: connection refused"),
			expCode:    http.StatusServiceUnavailable,
			expErrCode: "auth_unavailable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &authStackRoundTripper{authErr: test.authErr}
			hb := newStackHandlerBuilder(&http.Client{Transport: transport})

			h, err := hb.Build(config.Root{
				Server: stripUserIDServer(t),
				Routes: []config.Route{
					{
						Path:   "/chats/{chatID}/messages",
						Method: test.method,
						Checks: []config.Check{
							{Type: "policy", Config: mustYAML(t, map[string]any{
								"upstream": map[string]any{
									"scheme": "http",
									"host":   "auth.local",
									"path":   "/auth/validate",
									"method": http.MethodPost,
								},
								"expected_status": http.StatusOK,
								"store": map[string]any{
									"user_id": "{body:id}",
								},
							})},
						},
						Transforms: config.Transform{
							"header": mustYAML(t, map[string]any{
								"X-User-ID": "{context:user_id}",
							}),
						},
						Upstream: config.Upstream{
							Scheme: "http",
							Host:   "api.local",
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			req := httptest.NewRequest(test.method, test.target, strings.NewReader(test.body))
			req.Header.Set("X-User-ID", "forged-user")
			if test.token != "" {
				req.Header.Set("Authorization", test.token)
			}
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != test.expCode {
				t.Fatalf("expected status %d, got %d: %s", test.expCode, rec.Code, rec.Body.String())
			}

			if transport.authCalls != 1 || transport.authMethod != http.MethodPost || transport.authPath != "/auth/validate" {
				t.Fatalf("expected one POST /auth/validate, got %d calls, %s %s", transport.authCalls, transport.authMethod, transport.authPath)
			}
			if transport.authBody != "" {
				t.Fatalf("expected no client body in auth request, got %q", transport.authBody)
			}
			if transport.authHeader != test.token {
				t.Fatalf("expected Authorization %q in auth request, got %q", test.token, transport.authHeader)
			}
			if transport.authUserID != "" {
				t.Fatalf("expected no client X-User-ID in auth request, got %q", transport.authUserID)
			}

			if test.expErrCode != "" {
				if !strings.Contains(rec.Body.String(), `"code":"`+test.expErrCode+`"`) {
					t.Fatalf("expected error code %s, got %s", test.expErrCode, rec.Body.String())
				}
				if transport.apiCalls != 0 {
					t.Fatal("expected upstream not to be called")
				}
				return
			}

			if transport.apiUserID != "user-1" {
				t.Fatalf("expected verified X-User-ID, got %q", transport.apiUserID)
			}
			if transport.apiPath != test.expAPIPath || transport.apiQuery != test.expQuery {
				t.Fatalf("expected upstream %s?%s, got %s?%s", test.expAPIPath, test.expQuery, transport.apiPath, transport.apiQuery)
			}
			if transport.apiBody != test.body {
				t.Fatalf("expected upstream body %q, got %q", test.body, transport.apiBody)
			}
		})
	}
}

func stripUserIDServer(t *testing.T) config.Server {
	t.Helper()

	return config.Server{
		Middlewares: []config.Middleware{
			{Type: "strip_headers", Config: mustYAML(t, map[string]any{
				"headers": []string{"X-User-ID"},
			})},
		},
	}
}

// A route without checks: the header listed in strip_headers does not reach the service.
func TestStack_StripHeadersOnOpenRoute(t *testing.T) {
	transport := &authStackRoundTripper{}
	hb := newStackHandlerBuilder(&http.Client{Transport: transport})

	h, err := hb.Build(config.Root{
		Server: stripUserIDServer(t),
		Routes: []config.Route{
			{
				Path:   "/auth/login",
				Method: http.MethodPost,
				Upstream: config.Upstream{
					Scheme: "http",
					Host:   "api.local",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"alice"}`))
	req.Header.Set("X-User-ID", "forged-user")
	req.Header.Set("Authorization", "Bearer kept")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if transport.apiCalls != 1 || transport.authCalls != 0 {
		t.Fatalf("expected one upstream call and no auth calls, got %d and %d", transport.apiCalls, transport.authCalls)
	}
	if transport.apiUserID != "" {
		t.Fatalf("expected no X-User-ID in upstream request, got %q", transport.apiUserID)
	}
	if transport.apiAuthSeen != "Bearer kept" || transport.apiBody != `{"username":"alice"}` {
		t.Fatalf("expected other headers and body to stay, got %q and %q", transport.apiAuthSeen, transport.apiBody)
	}
}
