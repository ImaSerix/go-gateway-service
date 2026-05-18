package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/ctxkeys"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		allowedOrigins    []string
		allowedMethods    []string
		allowedHeaders    []string
		reqMethod         string
		reqOrigin         string
		expHandlerCalled  bool
		expCode           int
		expOrigin         string
		expMethods        string
		expAllowedHeaders string
		expVary           string
	}{
		{
			name: "allowed origin",
			allowedOrigins: []string{
				"http://allowed.origin",
				"http://allowed.second.origin",
			},
			allowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			allowedHeaders: []string{
				"X-Username",
				"X-Password",
			},
			reqMethod:         http.MethodGet,
			reqOrigin:         "http://allowed.origin",
			expHandlerCalled:  true,
			expCode:           http.StatusOK,
			expOrigin:         "http://allowed.origin",
			expMethods:        "GET, POST",
			expAllowedHeaders: "X-Username, X-Password",
			expVary:           "Origin",
		},
		{
			name: "options allowed origin",
			allowedOrigins: []string{
				"http://allowed.origin",
			},
			allowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			allowedHeaders: []string{
				"X-Username",
			},
			reqMethod:         http.MethodOptions,
			reqOrigin:         "http://allowed.origin",
			expHandlerCalled:  false,
			expCode:           http.StatusNoContent,
			expOrigin:         "http://allowed.origin",
			expMethods:        "GET, POST",
			expAllowedHeaders: "X-Username",
			expVary:           "Origin",
		},
		{
			name: "allow all origins",
			allowedOrigins: []string{
				"*",
			},
			allowedMethods: []string{
				http.MethodGet,
			},
			allowedHeaders: []string{
				"X-Token",
			},
			reqMethod:         http.MethodGet,
			reqOrigin:         "http://random.origin",
			expHandlerCalled:  true,
			expCode:           http.StatusOK,
			expOrigin:         "*",
			expMethods:        "GET",
			expAllowedHeaders: "X-Token",
		},
		{
			name: "empty origin",
			allowedOrigins: []string{
				"*",
			},
			allowedMethods: []string{
				http.MethodGet,
			},
			allowedHeaders: []string{
				"X-Token",
			},
			reqMethod:        http.MethodGet,
			expHandlerCalled: true,
			expCode:          http.StatusOK,
		},
		{
			name: "not allowed origin",
			allowedOrigins: []string{
				"http://allowed.origin",
			},
			allowedMethods: []string{
				http.MethodGet,
			},
			allowedHeaders: []string{
				"X-Token",
			},
			reqMethod:        http.MethodGet,
			reqOrigin:        "http://bad.origin",
			expHandlerCalled: true,
			expCode:          http.StatusOK,
		},
		{
			name:             "allowed origin without methods and headers",
			allowedOrigins:   []string{"http://allowed.origin"},
			reqMethod:        http.MethodGet,
			reqOrigin:        "http://allowed.origin",
			expHandlerCalled: true,
			expCode:          http.StatusOK,
			expOrigin:        "http://allowed.origin",
			expVary:          "Origin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var hcalled bool

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hcalled = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(test.reqMethod, "http://nice.url", nil)
			if test.reqOrigin != "" {
				req.Header.Set("Origin", test.reqOrigin)
			}
			w := httptest.NewRecorder()

			CORS(test.allowedOrigins, test.allowedMethods, test.allowedHeaders)(h).ServeHTTP(w, req)

			if hcalled != test.expHandlerCalled {
				t.Fatalf("expected next handler called=%v, but got %v", test.expHandlerCalled, hcalled)
			}
			if w.Code != test.expCode {
				t.Fatalf("expected status %d, but got %d", test.expCode, w.Code)
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
			if v := w.Header().Get("Vary"); v != test.expVary {
				t.Fatalf("expected Vary header %s, but got %s", test.expVary, v)
			}
		})
	}
}

func TestLoggingMiddleware_CallsNextHandler(t *testing.T) {
	var hcalled bool

	req := httptest.NewRequest(http.MethodGet, "http://nice.url", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hcalled = true
		w.WriteHeader(http.StatusCreated)
	})

	Logging(h).ServeHTTP(w, req)

	if !hcalled {
		t.Fatal("expected next handler called, but it was not")
	}
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestMetricMiddleware(t *testing.T) {
	m := NewMetric()
	var hcalled bool

	req := httptest.NewRequest(http.MethodPost, "http://nice.url/path", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hcalled = true
		w.WriteHeader(http.StatusCreated)
	})

	m.Middleware()(h).ServeHTTP(w, req)

	if !hcalled {
		t.Fatal("expected next handler called, but it was not")
	}
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	if got := m.request["POST /path"]; got != 1 {
		t.Fatalf("expected request count 1, but got %d", got)
	}
	if got := m.statusCode["POST /path 201"]; got != 1 {
		t.Fatalf("expected status count 1, but got %d", got)
	}
	if got := m.latency["POST /path"]; got <= 0 {
		t.Fatalf("expected positive latency, but got %v", got)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		expHandlerCalled bool
		expCode          int
	}{
		{
			name: "success",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			}),
			expHandlerCalled: true,
			expCode:          http.StatusAccepted,
		},
		{
			name: "panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("bad panic")
			}),
			expHandlerCalled: true,
			expCode:          http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var hcalled bool
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hcalled = true
				test.handler.ServeHTTP(w, r)
			})

			req := httptest.NewRequest(http.MethodGet, "http://nice.url", nil)
			w := httptest.NewRecorder()

			Recovery(h).ServeHTTP(w, req)

			if hcalled != test.expHandlerCalled {
				t.Fatalf("expected next handler called=%v, but got %v", test.expHandlerCalled, hcalled)
			}
			if w.Code != test.expCode {
				t.Fatalf("expected status %d, got %d", test.expCode, w.Code)
			}
		})
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		reqID       string
		expSameID   bool
		expValidLen bool
	}{
		{
			name:      "request id from header",
			reqID:     "request-id-123",
			expSameID: true,
		},
		{
			name:        "generated request id",
			expValidLen: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotID string
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				v, ok := r.Context().Value(ctxkeys.CtxRequestIDKey).(string)
				if !ok {
					t.Fatal("expected request id in context")
				}
				gotID = v
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "http://nice.url", nil)
			if test.reqID != "" {
				req.Header.Set("X-Request-ID", test.reqID)
			}
			w := httptest.NewRecorder()

			RequestID(h).ServeHTTP(w, req)

			if test.expSameID && gotID != test.reqID {
				t.Fatalf("expected request id %s, but got %s", test.reqID, gotID)
			}
			if test.expValidLen && len(gotID) != 36 {
				t.Fatalf("expected generated request id length 36, but got %d", len(gotID))
			}
			if v := w.Header().Get("X-Request-ID"); v != gotID {
				t.Fatalf("expected response request id %s, but got %s", gotID, v)
			}
		})
	}
}

func TestRealIPMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		remoteAddr   string
		forwardedFor string
		realIP       string
		expIP        string
	}{
		{
			name:       "remote addr",
			remoteAddr: "127.0.0.1:1234",
			expIP:      "127.0.0.1:1234",
		},
		{
			name:         "x forwarded for",
			remoteAddr:   "127.0.0.1:1234",
			forwardedFor: "10.0.0.1",
			expIP:        "10.0.0.1",
		},
		{
			name:         "x real ip has priority",
			remoteAddr:   "127.0.0.1:1234",
			forwardedFor: "10.0.0.1",
			realIP:       "10.0.0.2",
			expIP:        "10.0.0.2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotIP string
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				v, ok := r.Context().Value(ctxkeys.CtxRealIPKey).(string)
				if !ok {
					t.Fatal("expected real ip in context")
				}
				gotIP = v
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "http://nice.url", nil)
			req.RemoteAddr = test.remoteAddr
			if test.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", test.forwardedFor)
			}
			if test.realIP != "" {
				req.Header.Set("X-Real-IP", test.realIP)
			}
			w := httptest.NewRecorder()

			RealIP(h).ServeHTTP(w, req)

			if gotIP != test.expIP {
				t.Fatalf("expected real ip %s, but got %s", test.expIP, gotIP)
			}
		})
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	var hasDeadline bool
	var timeoutErr error

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasDeadline = r.Context().Deadline()
		<-r.Context().Done()
		timeoutErr = r.Context().Err()
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "http://nice.url", nil)
	w := httptest.NewRecorder()

	Timeout(10*time.Millisecond)(h).ServeHTTP(w, req)

	if !hasDeadline {
		t.Fatal("expected context with deadline")
	}
	if timeoutErr != context.DeadlineExceeded {
		t.Fatalf("expected context deadline exceeded, but got %v", timeoutErr)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		remote   string
		expCalls int
	}{
		{
			name:     "remote addr client id",
			ctx:      context.Background(),
			remote:   "127.0.0.1:1234",
			expCalls: 2,
		},
		{
			name:     "real ip client id",
			ctx:      context.WithValue(context.Background(), ctxkeys.CtxRealIPKey, "10.0.0.1"),
			remote:   "127.0.0.1:1234",
			expCalls: 2,
		},
		{
			name:     "user id client id",
			ctx:      context.WithValue(context.Background(), ctxkeys.CtxUserIDKey, "user-123"),
			remote:   "127.0.0.1:1234",
			expCalls: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rl := NewRateLimit(2, time.Second)
			var called int

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called++
				w.WriteHeader(http.StatusOK)
			})

			m := rl.Middleware()(h)
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil).WithContext(test.ctx)
			req.RemoteAddr = test.remote

			w1 := httptest.NewRecorder()
			m.ServeHTTP(w1, req)
			if w1.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", w1.Code)
			}

			w2 := httptest.NewRecorder()
			m.ServeHTTP(w2, req)
			if w2.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", w2.Code)
			}

			w3 := httptest.NewRecorder()
			m.ServeHTTP(w3, req)
			if w3.Code != http.StatusTooManyRequests {
				t.Fatalf("expected status 429, got %d", w3.Code)
			}

			if called != test.expCalls {
				t.Fatalf("expected handler called %d times, got %d", test.expCalls, called)
			}
		})
	}
}

func TestRateLimitMiddleware_ResetWindow(t *testing.T) {
	rl := NewRateLimit(1, 10*time.Millisecond)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m := rl.Middleware()(h)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	w1 := httptest.NewRecorder()
	m.ServeHTTP(w1, req)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	m.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", w2.Code)
	}

	time.Sleep(20 * time.Millisecond)

	w3 := httptest.NewRecorder()
	m.ServeHTTP(w3, req)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected status 200 after reset, got %d", w3.Code)
	}
}

func TestStatusWriter(t *testing.T) {
	t.Run("write header", func(t *testing.T) {
		w := httptest.NewRecorder()
		sw := &statusWriter{ResponseWriter: w}

		sw.WriteHeader(http.StatusCreated)

		if sw.statusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", sw.statusCode)
		}
		if w.Code != http.StatusCreated {
			t.Fatalf("expected response status 201, got %d", w.Code)
		}
	})

	t.Run("write default status", func(t *testing.T) {
		w := httptest.NewRecorder()
		sw := &statusWriter{ResponseWriter: w}

		if _, err := sw.Write([]byte("nice body")); err != nil {
			t.Fatalf("expected nil error, but got %v", err)
		}

		if sw.statusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", sw.statusCode)
		}
		if w.Code != http.StatusOK {
			t.Fatalf("expected response status 200, got %d", w.Code)
		}
		if body := w.Body.String(); body != "nice body" {
			t.Fatalf("expected body nice body, got %s", body)
		}
	})
}
