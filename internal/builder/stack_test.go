package builder_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
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
	authURI     string
	authBody    string
	authHeader  string
	authUserID  string
	apiCalls    int
	apiMethod   string
	apiPath     string
	apiQuery    string
	apiURI      string
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
		rt.authURI = r.URL.RequestURI()
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
		rt.apiMethod = r.Method
		rt.apiPath = r.URL.Path
		rt.apiQuery = r.URL.RawQuery
		rt.apiURI = r.URL.RequestURI()
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
			name:       "PUT without body",
			method:     http.MethodPut,
			target:     "/chats/7/messages",
			token:      "Bearer good",
			expCode:    http.StatusCreated,
			expAPIPath: "/chats/7/messages",
		},
		{
			name:       "DELETE without body",
			method:     http.MethodDelete,
			target:     "/chats/7/messages",
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
			if transport.apiMethod != test.method {
				t.Fatalf("expected upstream method %s, got %s", test.method, transport.apiMethod)
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

// upstream.path: a non-empty path replaces the client path, an empty one passes
// the client path through; the client query is always kept after the configured one.
// The request URI is compared as sent, so a lost path encoding is visible.
func TestStack_UpstreamPath(t *testing.T) {
	tests := []struct {
		name         string
		routePath    string
		upstreamPath string
		target       string

		expURI string
	}{
		{
			name:         "configured path equal to client path is not doubled",
			routePath:    "/resource/{id}",
			upstreamPath: "/resource/{route:id}",
			target:       "/resource/42?limit=5&before=10",
			expURI:       "/resource/42?limit=5&before=10",
		},
		{
			name:         "configured path replaces client path",
			routePath:    "/users/{id}",
			upstreamPath: "/api/v1/profiles/{route:id}",
			target:       "/users/42",
			expURI:       "/api/v1/profiles/42",
		},
		{
			name:      "empty path passes client path and query through",
			routePath: "/resource/{id}",
			target:    "/resource/42?limit=5&before=10",
			expURI:    "/resource/42?limit=5&before=10",
		},
		{
			name:         "configured query goes before client query",
			routePath:    "/find",
			upstreamPath: "/search?source=gateway&locale={query:locale}",
			target:       "/find?q=cats&locale=lv",
			expURI:       "/search?source=gateway&locale=lv&q=cats&locale=lv",
		},
		{
			name:         "configured query without client query",
			routePath:    "/find",
			upstreamPath: "/search?source=gateway",
			target:       "/find",
			expURI:       "/search?source=gateway",
		},
		{
			name:      "empty path passes root through",
			routePath: "/",
			target:    "/",
			expURI:    "/",
		},
		{
			name:         "configured root replaces client path",
			routePath:    "/health",
			upstreamPath: "/",
			target:       "/health",
			expURI:       "/",
		},
		{
			name:      "empty path keeps client trailing slash",
			routePath: "/resource/{id}/",
			target:    "/resource/42/",
			expURI:    "/resource/42/",
		},
		{
			name:         "configured trailing slash is kept",
			routePath:    "/resource/{id}",
			upstreamPath: "/resource/{route:id}/",
			target:       "/resource/42",
			expURI:       "/resource/42/",
		},
		{
			name:         "configured path without trailing slash drops client one",
			routePath:    "/resource/{id}/",
			upstreamPath: "/resource/{route:id}",
			target:       "/resource/42/",
			expURI:       "/resource/42",
		},
		{
			name:      "empty path keeps client path encoding",
			routePath: "/resource/{id}",
			target:    "/resource/a%2Fb?limit=5",
			expURI:    "/resource/a%2Fb?limit=5",
		},
		{
			name:         "configured path encoding is kept",
			routePath:    "/resource/{id}",
			upstreamPath: "/api/a%2Fb/{route:id}",
			target:       "/resource/42",
			expURI:       "/api/a%2Fb/42",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &authStackRoundTripper{}
			hb := newStackHandlerBuilder(&http.Client{Transport: transport})

			h, err := hb.Build(config.Root{
				Routes: []config.Route{
					{
						Path:   test.routePath,
						Method: http.MethodGet,
						Upstream: config.Upstream{
							Scheme: "http",
							Host:   "api.local",
							Path:   test.upstreamPath,
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, test.target, nil))

			if rec.Code != http.StatusCreated {
				t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
			}
			if transport.apiCalls != 1 {
				t.Fatalf("expected one upstream call, got %d", transport.apiCalls)
			}
			if transport.apiURI != test.expURI {
				t.Fatalf("expected upstream %s, got %s", test.expURI, transport.apiURI)
			}
		})
	}
}

// A template in upstream.path that cannot be resolved: 500 internal_error, no upstream request.
func TestStack_UpstreamPathTemplateError(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
	}))
	defer srv.Close()

	hb := newStackHandlerBuilder(srv.Client())

	h, err := hb.Build(config.Root{
		Routes: []config.Route{
			{
				Path:   "/resource",
				Method: http.MethodGet,
				Upstream: config.Upstream{
					Scheme: "http",
					Host:   srv.Listener.Addr().String(),
					Path:   "/resource/{query:id}",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/resource", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d: %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"code":"internal_error"`) {
		t.Fatalf("expected internal_error code, got %s", rec.Body.String())
	}
	if n := calls.Load(); n != 0 {
		t.Fatalf("expected upstream not to be called, got %d calls", n)
	}
}

// A value from the request stays a value: it cannot change the path or the query
// structure of the upstream URL, the path suffix from the config is kept.
func TestStack_UpstreamPathValuesAreEscaped(t *testing.T) {
	tests := []struct {
		name         string
		routePath    string
		upstreamPath string
		target       string

		expURI string
	}{
		{
			name:         "fragment sign in route value keeps path suffix",
			routePath:    "/users/{id}",
			upstreamPath: "/users/{route:id}/profile",
			target:       "/users/42%23",
			expURI:       "/users/42%23/profile",
		},
		{
			name:         "query and fragment signs in path value keep path suffix",
			routePath:    "/view",
			upstreamPath: "/public/{query:id}/info",
			target:       "/view?id=42%3Fadmin%3Dtrue%23",
			expURI:       "/public/42%3Fadmin=true%23/info?id=42%3Fadmin%3Dtrue%23",
		},
		{
			name:         "ampersand in query value does not add a parameter",
			routePath:    "/find",
			upstreamPath: "/search?locale={query:locale}",
			target:       "/find?locale=lv%26admin%3Dtrue",
			expURI:       "/search?locale=lv%26admin%3Dtrue&locale=lv%26admin%3Dtrue",
		},
		{
			name:         "slashes in path value do not add segments",
			routePath:    "/view",
			upstreamPath: "/public/{query:id}",
			target:       "/view?id=..%2Fadmin%2Fdelete",
			expURI:       "/public/..%2Fadmin%2Fdelete?id=..%2Fadmin%2Fdelete",
		},
		{
			name:         "encoded slash in route value is not encoded twice",
			routePath:    "/resource/{id}",
			upstreamPath: "/api/{route:id}/info",
			target:       "/resource/a%2Fb",
			expURI:       "/api/a%2Fb/info",
		},
		{
			name:         "percent sign in path value",
			routePath:    "/view",
			upstreamPath: "/public/{query:id}/info",
			target:       "/view?id=50%2541",
			expURI:       "/public/50%2541/info?id=50%2541",
		},
		{
			name:         "header value in path",
			routePath:    "/view",
			upstreamPath: "/public/{header:X-Name}/info",
			target:       "/view",
			expURI:       "/public/a%2Fb%3Fc/info",
		},
		{
			name:         "uuid is passed as before",
			routePath:    "/users/{id}",
			upstreamPath: "/api/users/{route:id}/profile?ref={query:ref}",
			target:       "/users/3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e?ref=Abc-123",
			expURI:       "/api/users/3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e/profile?ref=Abc-123&ref=Abc-123",
		},
		{
			name:         "dot segments written in config are kept",
			routePath:    "/users/{id}",
			upstreamPath: "/api/./v1/../users/{route:id}",
			target:       "/users/42",
			expURI:       "/api/./v1/../users/42",
		},
		{
			name:         "dots in query value are a regular value",
			routePath:    "/find",
			upstreamPath: "/search?dir={query:dir}",
			target:       "/find?dir=..",
			expURI:       "/search?dir=..&dir=..",
		},
		{
			name:         "value with dots that is not a dot segment",
			routePath:    "/files/{name}",
			upstreamPath: "/files/{route:name}",
			target:       "/files/..a",
			expURI:       "/files/..a",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &authStackRoundTripper{}
			hb := newStackHandlerBuilder(&http.Client{Transport: transport})

			h, err := hb.Build(config.Root{
				Routes: []config.Route{
					{
						Path:   test.routePath,
						Method: http.MethodGet,
						Upstream: config.Upstream{
							Scheme: "http",
							Host:   "api.local",
							Path:   test.upstreamPath,
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			req.Header.Set("X-Name", "a/b?c")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusCreated {
				t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
			}
			if transport.apiCalls != 1 {
				t.Fatalf("expected one upstream call, got %d", transport.apiCalls)
			}
			if transport.apiURI != test.expURI {
				t.Fatalf("expected upstream %s, got %s", test.expURI, transport.apiURI)
			}
		})
	}
}

// A value that turns a path segment into "." or "..": 400 bad_request, no upstream request.
func TestStack_UpstreamPathDotSegmentValue(t *testing.T) {
	tests := []struct {
		name         string
		routePath    string
		upstreamPath string
		target       string
	}{
		{name: "two dots from query", routePath: "/view", upstreamPath: "/public/{query:id}/info", target: "/view?id=.."},
		{name: "one dot from query", routePath: "/view", upstreamPath: "/public/{query:id}/info", target: "/view?id=."},
		{name: "two dots from route", routePath: "/users/{id}", upstreamPath: "/users/{route:id}/profile", target: "/users/.."},
		{name: "encoded dots from route", routePath: "/users/{id}", upstreamPath: "/users/{route:id}/profile", target: "/users/%2e%2e"},
		{name: "two dots from header", routePath: "/view", upstreamPath: "/public/{header:X-Name}", target: "/view"},
		{name: "dot segment made of two values", routePath: "/view", upstreamPath: "/public/{query:a}{query:b}/info", target: "/view?a=.&b=."},
		{name: "dot segment made of config and value", routePath: "/view", upstreamPath: "/public/.{query:a}/info", target: "/view?a=."},
		{name: "dot value next to config text", routePath: "/view", upstreamPath: "/public/v{query:a}/info", target: "/view?a=.."},
		{name: "one dot made of config and an empty value", routePath: "/view", upstreamPath: "/public/.{context:empty}/info", target: "/view"},
		{name: "two dots made of config and an empty value", routePath: "/view", upstreamPath: "/public/..{context:empty}/info", target: "/view"},
		{name: "dot value before an unresolved template", routePath: "/view", upstreamPath: "/public/{query:id}/{query:missing}", target: "/view?id=.."},
		{name: "dot value after an unresolved template", routePath: "/view", upstreamPath: "/public/{query:missing}/{query:id}", target: "/view?id=.."},
		{name: "dot value with an unresolved template in query", routePath: "/view", upstreamPath: "/public/{query:id}?x={query:missing}", target: "/view?id=.."},
		{name: "made dot segment before an unresolved template", routePath: "/view", upstreamPath: "/public/.{context:empty}/{query:missing}", target: "/view"},
		{name: "made dot segment after an unresolved template", routePath: "/view", upstreamPath: "/public/{query:missing}/.{context:empty}", target: "/view"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
			}))
			defer srv.Close()

			hb := newStackHandlerBuilder(srv.Client())

			h, err := hb.Build(config.Root{
				Routes: []config.Route{
					{
						Path:   test.routePath,
						Method: http.MethodGet,
						Middlewares: []config.Middleware{
							{Type: "inject", Config: mustYAML(t, map[string]any{
								"context": map[string]any{"empty": ""},
							})},
						},
						Upstream: config.Upstream{
							Scheme: "http",
							Host:   srv.Listener.Addr().String(),
							Path:   test.upstreamPath,
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			req.Header.Set("X-Name", "..")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("expected json content type, got %q", got)
			}
			if !strings.Contains(rec.Body.String(), `"code":"bad_request"`) {
				t.Fatalf("expected bad_request code, got %s", rec.Body.String())
			}
			if strings.Contains(rec.Body.String(), "..") {
				t.Fatalf("value leaked to client: %s", rec.Body.String())
			}
			if n := calls.Load(); n != 0 {
				t.Fatalf("expected upstream not to be called, got %d calls", n)
			}
		})
	}
}

// The same rules for the URL of a policy request.
func TestStack_PolicyPathValues(t *testing.T) {
	tests := []struct {
		name       string
		policyPath string
		target     string

		expStatus  int
		expCode    string
		expAuthURI string
	}{
		{
			name:       "fragment sign in route value keeps path suffix",
			policyPath: "/auth/{route:id}/check",
			target:     "/resource/42%23",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/42%23/check",
		},
		{
			name:       "query and fragment signs in path value keep path suffix",
			policyPath: "/auth/{query:id}/check",
			target:     "/resource/1?id=42%3Fadmin%3Dtrue%23",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/42%3Fadmin=true%23/check",
		},
		{
			name:       "ampersand in query value does not add a parameter",
			policyPath: "/auth/check?locale={query:locale}",
			target:     "/resource/1?locale=lv%26admin%3Dtrue",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/check?locale=lv%26admin%3Dtrue",
		},
		{
			name:       "slashes in path value do not add segments",
			policyPath: "/auth/{query:id}",
			target:     "/resource/1?id=..%2Fadmin%2Fdelete",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/..%2Fadmin%2Fdelete",
		},
		{
			name:       "uuid is passed as before",
			policyPath: "/auth/{route:id}/check",
			target:     "/resource/3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/3f2b8c1e-7a4d-4e9b-9c55-0d1a2b3c4d5e/check",
		},
		{
			name:       "two dots in path value",
			policyPath: "/auth/{query:id}/check",
			target:     "/resource/1?id=..",
			expStatus:  http.StatusBadRequest,
			expCode:    "bad_request",
		},
		{
			name:       "one dot in route value",
			policyPath: "/auth/{route:id}/check",
			target:     "/resource/.",
			expStatus:  http.StatusBadRequest,
			expCode:    "bad_request",
		},
		{
			name:       "encoded dots in route value",
			policyPath: "/auth/{route:id}/check",
			target:     "/resource/%2e%2e",
			expStatus:  http.StatusBadRequest,
			expCode:    "bad_request",
		},
		{
			name:       "encoded slash in route value is not encoded twice",
			policyPath: "/auth/{route:id}/check",
			target:     "/resource/a%2Fb",
			expStatus:  http.StatusCreated,
			expAuthURI: "/auth/a%2Fb/check",
		},
		{
			name:       "dot value before an unresolved template",
			policyPath: "/auth/{query:id}/{query:missing}",
			target:     "/resource/1?id=..",
			expStatus:  http.StatusBadRequest,
			expCode:    "bad_request",
		},
		{
			name:       "dot value after an unresolved template",
			policyPath: "/auth/{query:missing}/{query:id}",
			target:     "/resource/1?id=..",
			expStatus:  http.StatusBadRequest,
			expCode:    "bad_request",
		},
		{
			name:       "unresolved template without a dot value",
			policyPath: "/auth/{query:missing}/check",
			target:     "/resource/1",
			expStatus:  http.StatusServiceUnavailable,
			expCode:    "auth_unavailable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &authStackRoundTripper{}
			hb := newStackHandlerBuilder(&http.Client{Transport: transport})

			h, err := hb.Build(config.Root{
				Routes: []config.Route{
					{
						Path:   "/resource/{id}",
						Method: http.MethodGet,
						Checks: []config.Check{
							{Type: "policy", Config: mustYAML(t, map[string]any{
								"upstream": map[string]any{
									"scheme": "http",
									"host":   "auth.local",
									"path":   test.policyPath,
									"method": http.MethodPost,
								},
								"expected_status": http.StatusOK,
								"store": map[string]any{
									"user_id": "{body:id}",
								},
							})},
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

			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			req.Header.Set("Authorization", "Bearer good")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != test.expStatus {
				t.Fatalf("expected status %d, got %d: %s", test.expStatus, rec.Code, rec.Body.String())
			}

			if test.expCode != "" {
				if !strings.Contains(rec.Body.String(), `"code":"`+test.expCode+`"`) {
					t.Fatalf("expected %s code, got %s", test.expCode, rec.Body.String())
				}
				if transport.authCalls != 0 || transport.apiCalls != 0 {
					t.Fatalf("expected no requests, got %d policy and %d upstream calls", transport.authCalls, transport.apiCalls)
				}
				return
			}

			if transport.authCalls != 1 || transport.apiCalls != 1 {
				t.Fatalf("expected one policy and one upstream call, got %d and %d", transport.authCalls, transport.apiCalls)
			}
			if transport.authURI != test.expAuthURI {
				t.Fatalf("expected policy request %s, got %s", test.expAuthURI, transport.authURI)
			}
		})
	}
}

// Escaping belongs to URLs only: transforms get a route value exactly as the router gives it,
// for an escaped request path that is the escaped segment.
func TestStack_RouteValueInTransformsIsNotChanged(t *testing.T) {
	transport := &authStackRoundTripper{}
	hb := newStackHandlerBuilder(&http.Client{Transport: transport})

	h, err := hb.Build(config.Root{
		Routes: []config.Route{
			{
				Path:   "/r/{id}",
				Method: http.MethodPost,
				Transforms: config.Transform{
					"header":       mustYAML(t, map[string]any{"X-User-ID": "{route:id}"}),
					"query_params": mustYAML(t, map[string]any{"copy": "{route:id}"}),
					"body_fields":  mustYAML(t, map[string]any{"id": "{route:id}"}),
				},
				Upstream: config.Upstream{
					Scheme: "http",
					Host:   "api.local",
					Path:   "/api/{route:id}",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/r/a%2Fb", strings.NewReader(`{"name":"Ada"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if transport.apiUserID != "a%2Fb" {
		t.Fatalf("expected header value a%%2Fb, got %q", transport.apiUserID)
	}
	if !strings.Contains(transport.apiBody, `"id":"a%2Fb"`) {
		t.Fatalf("expected body field a%%2Fb, got %s", transport.apiBody)
	}
	// The path of the same request has the value encoded once.
	if transport.apiURI != "/api/a%2Fb?copy=a%252Fb" {
		t.Fatalf("expected upstream /api/a%%2Fb?copy=a%%252Fb, got %s", transport.apiURI)
	}
}
