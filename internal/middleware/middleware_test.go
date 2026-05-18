package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	buildermiddleware "github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/middleware"
	"gopkg.in/yaml.v3"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		reqMethod         string
		reqOrigin         string
		cfg               config.CORSMiddleaware
		expHandlerCalled  bool
		expOrigin         string
		expMethods        string
		expAllowedHeaders string
		expVary           bool
		expCode           int
	}{
		{
			name:      "success",
			reqMethod: "GET",
			reqOrigin: "http://allowed.origin",
			cfg: config.CORSMiddleaware{
				Allowed: config.CORSMiddleaware_Allowed{
					Origin: []string{
						"http://allowed.origin",
						"http://allowed.second.origin",
					},
					Method: []string{
						"GET",
						"POST",
					},
					Header: []string{
						"X-Username",
						"X-Password",
					},
				},
			},
			expVary:           true,
			expHandlerCalled:  true,
			expOrigin:         "http://allowed.origin",
			expMethods:        "GET, POST",
			expAllowedHeaders: "X-Username, X-Password",
			expCode:           200,
		},
		{
			name:      "success OPTIONS request",
			reqMethod: "OPTIONS",
			reqOrigin: "http://allowed.origin",
			cfg: config.CORSMiddleaware{
				Allowed: config.CORSMiddleaware_Allowed{
					Origin: []string{
						"http://allowed.origin",
						"http://allowed.second.origin",
					},
					Method: []string{
						"GET",
						"POST",
					},
					Header: []string{
						"X-Username",
						"X-Password",
					},
				},
			},
			expVary:           true,
			expHandlerCalled:  false,
			expOrigin:         "http://allowed.origin",
			expMethods:        "GET, POST",
			expAllowedHeaders: "X-Username, X-Password",
			expCode:           204,
		},
		{
			name:      "origin *",
			reqMethod: "GET",
			reqOrigin: "http://allowed.origin",
			cfg: config.CORSMiddleaware{
				Allowed: config.CORSMiddleaware_Allowed{
					Origin: []string{
						"*",
					},
					Method: []string{
						"GET",
						"POST",
					},
					Header: []string{
						"X-Username",
						"X-Password",
					},
				},
			},
			expHandlerCalled:  true,
			expOrigin:         "*",
			expMethods:        "GET, POST",
			expAllowedHeaders: "X-Username, X-Password",
			expCode:           204,
		},
		{
			name:      "request origin ''",
			reqMethod: "GET",
			reqOrigin: "",
			cfg: config.CORSMiddleaware{
				Allowed: config.CORSMiddleaware_Allowed{
					Origin: []string{
						"*",
					},
					Method: []string{
						"GET",
						"POST",
					},
					Header: []string{
						"X-Username",
						"X-Password",
					},
				},
			},
			expHandlerCalled:  true,
			expOrigin:         "",
			expMethods:        "",
			expAllowedHeaders: "",
			expCode:           200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var raw yaml.Node

			if err := raw.Encode(test.cfg); err != nil {
				t.Fatal(err)
			}

			m, err := buildermiddleware.NewCorsFactory().Create(raw)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}

			req := httptest.NewRequest(test.reqMethod, "http://nice.url", nil)
			w := httptest.NewRecorder()

			hcalled := false

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hcalled = true
				w.WriteHeader(http.StatusOK)
			})

			req.Header.Set("Origin", test.reqOrigin)

			m(h).ServeHTTP(w, req)

			if hcalled != test.expHandlerCalled {
				t.Fatalf("expected next handler called=%v, but got %v", test.expHandlerCalled, hcalled)
			}

			if v := w.Header().Get("Vary"); test.expVary && v != "Origin" {
				t.Fatal("expected response with 'Vary' header")
			}

			if v := w.Header().Get("Access-Control-Allow-Origin"); v != test.expOrigin {
				t.Fatalf("expected origin %s, but got %s", test.expOrigin, v)
			}

			if v := w.Header().Get("Access-Control-Allow-Methods"); v != test.expMethods {
				t.Fatalf("expected methods %s, but got %s", test.expMethods, v)
			}

			if v := w.Header().Get("Access-Control-Allow-Headers"); v != test.expAllowedHeaders {
				t.Fatalf("expected headers %s, but got %s", test.expAllowedHeaders, v)
			}

		})
	}
}

func TestLoggingMiddleware(t *testing.T) {

	var hcalled bool

	req := httptest.NewRequest("GET", "http://nice.url", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hcalled = true
		w.WriteHeader(http.StatusOK)
	})

	var raw yaml.Node

	m, err := buildermiddleware.NewLoggingFactory().Create(raw)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	m(h).ServeHTTP(w, req)

	if !hcalled {
		t.Fatal("expected next handler called, but it was not")
	}

}

func TestMetricMiddleware(t *testing.T) {

	var hcalled bool

	req := httptest.NewRequest("GET", "http://nice.url", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hcalled = true
		w.WriteHeader(http.StatusOK)
	})

	var raw yaml.Node

	m, err := buildermiddleware.NewMetricFactory().Create(raw)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	m(h).ServeHTTP(w, req)

	if !hcalled {
		t.Fatal("expected next handler called, but it was not")
	}

}

func TestRateLimitMiddleware(t *testing.T) {

	rl := middleware.NewRateLimit(2, time.Second)

	var called int

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	})

	m := rl.Middleware()(h)

	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	req.RemoteAddr = "127.0.0.1:1234"

	// first request
	w1 := httptest.NewRecorder()
	m.ServeHTTP(w1, req)

	if w1.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			w1.Code,
		)
	}

	// second request
	w2 := httptest.NewRecorder()
	m.ServeHTTP(w2, req)

	if w2.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			w2.Code,
		)
	}

	// third request -> should be limited
	w3 := httptest.NewRecorder()
	m.ServeHTTP(w3, req)

	if w3.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected status 429, got %d",
			w3.Code,
		)
	}

	if called != 2 {
		t.Fatalf(
			"expected handler called 2 times, got %d",
			called,
		)
	}
}

func TestRateLimitMiddleware_ResetWindow(t *testing.T) {

	rl := middleware.NewRateLimit(
		1,
		100*time.Millisecond,
	)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m := rl.Middleware()(h)

	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	req.RemoteAddr = "127.0.0.1:1234"

	// first request
	w1 := httptest.NewRecorder()
	m.ServeHTTP(w1, req)

	// second request -> limited
	w2 := httptest.NewRecorder()
	m.ServeHTTP(w2, req)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected status 429, got %d",
			w2.Code,
		)
	}

	time.Sleep(150 * time.Millisecond)

	// should work again
	w3 := httptest.NewRecorder()
	m.ServeHTTP(w3, req)

	if w3.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200 after reset, got %d",
			w3.Code,
		)
	}
}
