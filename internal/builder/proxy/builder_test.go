package proxy_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	proxyBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type rendererMock struct{}

func (rm *rendererMock) Render(s string, r *http.Request) (string, error) {
	return "", nil
}

func (rm *rendererMock) RenderURL(s string, r *http.Request) (string, error) {
	return rm.Render(s, r)
}

func TestProxyBuilder(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.Upstream
		expErr error
	}{
		{
			name: "succes",
			cfg: config.Upstream{
				Host:   "nice.host",
				Scheme: "http",
				Path:   "",
				Method: "GET",
			},
			expErr: nil,
		},
		{
			name: "new url error",
			cfg: config.Upstream{
				Host:   "",
				Scheme: "http",
				Method: "GET",
			},
			expErr: proxy.ErrInvalidHost,
		},
		{
			name: "invalid method",
			cfg: config.Upstream{
				Host:   "nice.host",
				Scheme: "http",
				Path:   "",
				Method: "INVALID",
			},
			expErr: proxy.ErrInvalidMethod,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			b := proxyBuilder.NewBuilder(http.DefaultClient, &rendererMock{})

			p, err := b.Build(test.cfg)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, but got %v", test.expErr, err)
			}

			if err == nil && p == nil {
				t.Fatal("got no error, but proxy is nil")
			}
			if err != nil && p != nil {
				t.Fatal("got no error, but proxy is not nil")
			}

		})
	}
}

type failingTransport struct{}

func (failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("dial tcp 10.0.0.1:8080: connection refused")
}

func TestProxyBuilder_UpstreamUnavailable(t *testing.T) {
	b := proxyBuilder.NewBuilder(&http.Client{Transport: failingTransport{}}, &rendererMock{})

	p, err := b.Build(config.Upstream{Host: "nice.host", Scheme: "http", Method: "GET"})
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	w := httptest.NewRecorder()
	p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users", nil))

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, but got %d", http.StatusBadGateway, w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"upstream_unavailable"`) {
		t.Fatalf("expected upstream_unavailable code, but got %s", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "10.0.0.1") {
		t.Fatalf("internal error text leaked to client: %s", w.Body.String())
	}
}

type fixedRenderer struct {
	out string
	err error
}

func (fr *fixedRenderer) Render(s string, r *http.Request) (string, error) {
	return fr.out, fr.err
}

func (fr *fixedRenderer) RenderURL(s string, r *http.Request) (string, error) {
	return fr.Render(s, r)
}

// captureLog collects the gateway log of one test.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	buf := &bytes.Buffer{}
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	return buf
}

// A failed upstream URL is a gateway error: 500 internal_error, no upstream request.
// A value that makes a dot segment is a client error: 400 bad_request, no upstream request.
// The answer does not depend on the state of the request context, the reason goes to the log.
func TestProxyBuilder_UpstreamURLError(t *testing.T) {
	renders := []struct {
		name   string
		render *fixedRenderer

		expStatus int
		expCode   string
		expLog    string
	}{
		{
			name:      "template render error",
			render:    &fixedRenderer{err: errors.New("failed resolve: secret-detail")},
			expStatus: http.StatusInternalServerError,
			expCode:   "internal_error",
			expLog:    "upstream url build failed",
		},
		{
			name:      "rendered url parse error",
			render:    &fixedRenderer{out: "http://secret-detail:port/users"},
			expStatus: http.StatusInternalServerError,
			expCode:   "internal_error",
			expLog:    "upstream url build failed",
		},
		{
			name:      "dot segment value",
			render:    &fixedRenderer{err: fmt.Errorf("%w: secret-detail", renderer.ErrDotSegment)},
			expStatus: http.StatusBadRequest,
			expCode:   "bad_request",
			expLog:    "upstream url rejected",
		},
	}

	contexts := []struct {
		name string
		ctx  func() (context.Context, context.CancelFunc)
	}{
		{
			name: "active context",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
		},
		{
			name: "cancelled context",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
		},
		{
			name: "expired deadline",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
			},
		},
	}

	for _, rend := range renders {
		for _, c := range contexts {
			t.Run(rend.name+"/"+c.name, func(t *testing.T) {
				var calls atomic.Int32
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
				}))
				defer srv.Close()

				log := captureLog(t)

				b := proxyBuilder.NewBuilder(srv.Client(), rend.render)

				p, err := b.Build(config.Upstream{
					Host:   srv.Listener.Addr().String(),
					Scheme: "http",
					Path:   "/users/{route:id}",
					Method: "GET",
				})
				if err != nil {
					t.Fatalf("expected nil error, but got %v", err)
				}

				ctx, cancel := c.ctx()
				defer cancel()

				// The incoming URL is absolute and points to the upstream itself:
				// the request must not be sent even to the address the client asked for.
				w := httptest.NewRecorder()
				p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, srv.URL+"/users/42", nil).WithContext(ctx))

				if w.Code != rend.expStatus {
					t.Fatalf("expected status %d, but got %d", rend.expStatus, w.Code)
				}
				if !strings.Contains(w.Body.String(), `"code":"`+rend.expCode+`"`) {
					t.Fatalf("expected %s code, but got %s", rend.expCode, w.Body.String())
				}
				if strings.Contains(w.Body.String(), "secret-detail") {
					t.Fatalf("internal error text leaked to client: %s", w.Body.String())
				}
				if n := calls.Load(); n != 0 {
					t.Fatalf("expected upstream not to be called, but got %d calls", n)
				}
				// The reason is kept in the log whatever the state of the context is.
				if !strings.Contains(log.String(), rend.expLog) || !strings.Contains(log.String(), "secret-detail") {
					t.Fatalf("expected %q with the reason in the log, but got %s", rend.expLog, log.String())
				}
			})
		}
	}
}

// A request cancelled by the client or by a deadline, with a valid upstream URL,
// is a transport error: 502 upstream_unavailable, not a gateway 500.
func TestProxyBuilder_CancelledRequest(t *testing.T) {
	contexts := []struct {
		name string
		ctx  func() (context.Context, context.CancelFunc)
	}{
		{
			name: "cancelled context",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
		},
		{
			name: "expired deadline",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
			},
		},
	}

	for _, c := range contexts {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			defer srv.Close()

			log := captureLog(t)

			b := proxyBuilder.NewBuilder(srv.Client(), &fixedRenderer{out: srv.URL + "/users/42"})

			p, err := b.Build(config.Upstream{
				Host:   srv.Listener.Addr().String(),
				Scheme: "http",
				Path:   "/users/{route:id}",
				Method: "GET",
			})
			if err != nil {
				t.Fatalf("expected nil error, but got %v", err)
			}

			ctx, cancel := c.ctx()
			defer cancel()

			w := httptest.NewRecorder()
			p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/42", nil).WithContext(ctx))

			if w.Code != http.StatusBadGateway {
				t.Fatalf("expected status %d, but got %d", http.StatusBadGateway, w.Code)
			}
			if !strings.Contains(w.Body.String(), `"code":"upstream_unavailable"`) {
				t.Fatalf("expected upstream_unavailable code, but got %s", w.Body.String())
			}
			if !strings.Contains(log.String(), "upstream request failed") || strings.Contains(log.String(), "upstream url build failed") {
				t.Fatalf("expected a transport error in the log, but got %s", log.String())
			}
		})
	}
}
