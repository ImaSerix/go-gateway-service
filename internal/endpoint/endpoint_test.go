package endpoint_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
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

func TestEndpoint_ServeHTTP(t *testing.T) {

	tests := []struct {
		name         string
		reqMethod    string
		endpointMeth string

		checkErr       error
		transformErr   error
		middlewarePass bool

		expCode    int
		expErrCode string
		expCalls   []string
	}{
		{
			name:         "success",
			reqMethod:    "GET",
			endpointMeth: "GET",

			middlewarePass: true,

			expCode: http.StatusOK,
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
				"check-2",
				"transform-1",
				"transform-2",
				"proxy",
			},
		},
		{
			name:         "method not allowed",
			reqMethod:    "POST",
			endpointMeth: "GET",

			expCode:    http.StatusMethodNotAllowed,
			expErrCode: "method_not_allowed",
			expCalls:   nil,
		},
		{
			name:         "middleware failed",
			reqMethod:    "GET",
			endpointMeth: "GET",

			middlewarePass: false,

			expCode: http.StatusInternalServerError,
			expCalls: []string{
				"middleware-1",
			},
		},
		{
			name:         "first check failed",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       check.ErrForbidden,
			middlewarePass: true,

			expCode:    http.StatusForbidden,
			expErrCode: "forbidden",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "check failed: no header",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       fmt.Errorf("%w: X-Request-ID", check.ErrNoHeader),
			middlewarePass: true,

			expCode:    http.StatusForbidden,
			expErrCode: "forbidden",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "check failed: bad request",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       fmt.Errorf("%w: value makes a dot segment", check.ErrBadRequest),
			middlewarePass: true,

			expCode:    http.StatusBadRequest,
			expErrCode: "bad_request",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "check failed: unauthorized",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       check.ErrUnauthorized,
			middlewarePass: true,

			expCode:    http.StatusUnauthorized,
			expErrCode: "unauthorized",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "check failed: policy unavailable",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       fmt.Errorf("%w: dial tcp 10.0.0.1:8080", check.ErrPolicyUnavailable),
			middlewarePass: true,

			expCode:    http.StatusServiceUnavailable,
			expErrCode: "auth_unavailable",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "check failed: internal error",
			reqMethod:    "GET",
			endpointMeth: "GET",

			checkErr:       http.ErrAbortHandler,
			middlewarePass: true,

			expCode:    http.StatusInternalServerError,
			expErrCode: "internal_error",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
			},
		},
		{
			name:         "transform failed",
			reqMethod:    "GET",
			endpointMeth: "GET",

			transformErr:   http.ErrBodyNotAllowed,
			middlewarePass: true,

			expCode:    http.StatusInternalServerError,
			expErrCode: "internal_error",
			expCalls: []string{
				"middleware-1",
				"middleware-2",
				"check-1",
				"check-2",
				"transform-1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var calls []string

			check1 := &mockCheck{
				name:   "check-1",
				calls:  &calls,
				err:    test.checkErr,
				setCtx: true,
			}

			check2 := &mockCheck{
				name:  "check-2",
				calls: &calls,
			}

			transform1 := &mockTransformer{
				name:      "transform-1",
				calls:     &calls,
				err:       test.transformErr,
				expectCtx: true,
			}

			transform2 := &mockTransformer{
				name:  "transform-2",
				calls: &calls,
			}

			proxy := &mockProxy{
				calls: &calls,
			}

			middleware1 := &mockMiddleware{
				name:  "middleware-1",
				pass:  test.middlewarePass,
				calls: &calls,
			}

			middleware2 := &mockMiddleware{
				name:  "middleware-2",
				pass:  true,
				calls: &calls,
			}

			path, _ := endpoint.NewPath("/users")
			method, _ := endpoint.NewMethod(test.endpointMeth)

			e := endpoint.NewEndpoint(
				path,
				method,
				[]pipeline.Checker{
					check1,
					check2,
				},
				[]pipeline.Transformer{
					transform1,
					transform2,
				},
				proxy,
				[]pipeline.Middleware{
					middleware1.Middleware(),
					middleware2.Middleware(),
				},
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.reqMethod, "http://nice.url", nil)

			e.ServeHTTP(w, r)

			if w.Code != test.expCode {
				t.Fatalf(
					"expected status code %d, but got %d",
					test.expCode,
					w.Code,
				)
			}

			if test.expErrCode != "" {
				var body struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatalf("expected json error body, got %q", w.Body.String())
				}
				if body.Code != test.expErrCode || body.Message == "" {
					t.Fatalf("expected error code %q with message, got %+v", test.expErrCode, body)
				}
				if ct := w.Header().Get("Content-Type"); ct != "application/json" {
					t.Fatalf("expected Content-Type application/json, got %q", ct)
				}
				if strings.Contains(w.Body.String(), "10.0.0.1") {
					t.Fatalf("internal error text leaked to client: %q", w.Body.String())
				}
			}

			if !reflect.DeepEqual(calls, test.expCalls) {
				t.Fatalf(
					"expected calls %v, but got %v",
					test.expCalls,
					calls,
				)
			}
		})
	}
}
