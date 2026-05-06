package endpoint_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
)

func TestPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expPath  string
		expError error
	}{
		{name: "success", path: "/nice/path", expPath: "/nice/path", expError: nil},
		{name: "success spaces", path: "   /nice/path    ", expPath: "/nice/path", expError: nil},
		{name: "empty path", path: "", expError: endpoint.ErrEmptyPath},
		{name: "path with spaces", path: "/bad path", expError: endpoint.ErrPathHasSpaces},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := endpoint.NewPath(test.path)
			if !errors.Is(err, test.expError) {
				t.Fatalf("expected wrapped error %v, but got %v", test.expError, err)
			}
			if err == nil && string(p) != test.expPath {
				t.Fatalf("expected path %s, but got %s", test.expPath, p)
			}
			if err != nil && string(p) != "" {
				t.Fatalf("expected empty path, but got %s", p)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expUrl   string
		expError error
	}{
		{name: "success", url: "http://nice.url", expUrl: "http://nice.url", expError: nil},
		{name: "unsupported scheme", url: "us://nice.url", expError: endpoint.ErrUnsupportedScheme},
		{name: "empty host", url: "http://", expError: endpoint.ErrEmptyHost},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := endpoint.NewURL(test.url)
			if !errors.Is(err, test.expError) {
				t.Fatalf("expected wrapped error %v, but got %v", test.expError, err)
			}
			if err == nil && string(u) != test.expUrl {
				t.Fatalf("expected path %s, but got %s", test.expUrl, u)
			}
			if err != nil && string(u) != "" {
				t.Fatalf("expected empty path, but got %s", u)
			}
		})
	}
}

func TestEndpoint(t *testing.T) {
	p, _ := endpoint.NewPath("/proxy/ping")
	url, _ := endpoint.NewURL("http://localhost:8080/ping")
	e := endpoint.NewEndpoint(p, endpoint.GET, url, endpoint.GET)

	if e == nil {
		t.Fatal("expected non-nil enpoint, but got nil")
	}
	if e.Path != p {
		t.Fatalf("expected path %s, but got %s", p, e.Path)
	}
	if e.Method != endpoint.GET {
		t.Fatalf("expected method %s, but got %s", endpoint.GET, e.Method)
	}
	if e.Upstream == nil {
		t.Fatalf("expected non-nil upstream, but got nil")
	}

	if e.Upstream.URL != url {
		t.Fatalf("expected target %s, but got %s", url, e.Upstream.URL)
	}
	if e.Upstream.Method != endpoint.GET {
		t.Fatalf("expected upstream method %s, but got %s", endpoint.GET, e.Upstream.Method)
	}
}

func TestEndpointFromConfig(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.RouteConfig
		expError error
	}{
		{
			name: "success",
			cfg: &config.RouteConfig{
				Path:   "/nice/path",
				Method: "GET",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "GET",
				},
			},
			expError: nil,
		},
		{
			name:     "nil config",
			cfg:      nil,
			expError: endpoint.ErrInvalidConfig,
		},
		{
			name: "invalid path",
			cfg: &config.RouteConfig{
				Path:   "",
				Method: "GET",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "GET",
				},
			},
			expError: endpoint.ErrEmptyPath,
		},
		{
			name: "invalid method",
			cfg: &config.RouteConfig{
				Path:   "/nice/path",
				Method: "INVALID",
				Upstream: config.Upstream{
					URL:    "https://jsonplaceholder.typicode.com/users",
					Method: "GET",
				},
			},
			expError: endpoint.ErrInvalidMethod,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			end, err := endpoint.NewEndpointFromConfig(test.cfg)
			if err != test.expError {
				t.Fatalf("expected error %v, but got %v", test.expError, err)
			}
			if end != nil && err != nil {
				t.Fatal("got erorr, but endpoint is not nil")
			}
			if end == nil && err == nil {
				t.Fatal("got nil endpoint, but no error")
			}
			if end != nil && test.cfg != nil {
				if end.Path != endpoint.Path(test.cfg.Path) {
					t.Fatalf("expected path %s, but got %s", test.cfg.Path, end.Path)
				}
				if end.Method != endpoint.Method(test.cfg.Method) {
					t.Fatalf("expected method %s, but got %s", test.cfg.Method, end.Method)
				}
				if end.Upstream == nil {
					t.Fatal("expected non-nil upstream, but got nil")
				}
				if end.Upstream.URL != endpoint.URL(test.cfg.Upstream.URL) {
					t.Fatalf("expected upstream url %s, but got %s", test.cfg.Upstream.URL, end.Upstream.URL)
				}
				if end.Upstream.Method != endpoint.Method(test.cfg.Upstream.Method) {
					t.Fatalf("expected upstream method %s, but got %s", test.cfg.Upstream.Method, end.Upstream.Method)
				}
			}
		})
	}
}

func TestEndpoint_ServeHTTP(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("POST /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong pong"))
	})
	mux.HandleFunc("GET /bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	})

	upstream := httptest.NewServer(mux)
	defer upstream.Close()

	path, _ := endpoint.NewPath("/proxy/ping")
	tURL, _ := endpoint.NewURL(upstream.URL + "/ping")
	tURLBad, _ := endpoint.NewURL(upstream.URL + "/bad")
	tBrokenURL, _ := endpoint.NewURL("http://notexisting.host/")

	tests := []struct {
		name         string
		method       endpoint.Method
		targetMethod endpoint.Method
		targetURL    endpoint.URL
		reqMethod    string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "GET -> GET upstream",
			method:       endpoint.GET,
			targetMethod: endpoint.GET,
			targetURL:    tURL,
			reqMethod:    "GET",
			expectedCode: http.StatusOK,
			expectedBody: "pong",
		},
		{
			name:         "GET -> POST upstream",
			method:       endpoint.GET,
			targetMethod: endpoint.POST,
			targetURL:    tURL,
			reqMethod:    "GET",
			expectedCode: http.StatusOK,
			expectedBody: "pong pong",
		},
		{
			name:         "wrong method",
			method:       endpoint.GET,
			targetMethod: endpoint.GET,
			targetURL:    tURL,
			reqMethod:    "POST",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "method not allowed\n",
		},
		{
			name:         "GET -> GET upstream (status 400)",
			method:       endpoint.GET,
			targetMethod: endpoint.GET,
			targetURL:    tURLBad,
			reqMethod:    "GET",
			expectedCode: http.StatusBadRequest,
			expectedBody: "bad request",
		},
		{
			name:         "upstream unavailable",
			method:       endpoint.GET,
			targetMethod: endpoint.GET,
			targetURL:    tBrokenURL,
			reqMethod:    "GET",
			expectedCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := endpoint.NewEndpoint(path, tt.method, tt.targetURL, tt.targetMethod)

			req := httptest.NewRequest(
				tt.reqMethod,
				"http://localhost:8080"+string(path),
				nil,
			)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, rec.Code)
			}

			body, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Fatalf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}
