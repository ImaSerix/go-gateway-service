package endpoint_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	endpointBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type checkerMock struct{}
func (c checkerMock) Execute(ctx context.Context, r *http.Request) (context.Context, error) { return ctx, nil }

type checkBuilderMock struct { many []pipeline.Checker; err error }
func (m checkBuilderMock) Build(config.Check) (pipeline.Checker, error) { return checkerMock{}, nil }
func (m checkBuilderMock) BuildMany([]config.Check) ([]pipeline.Checker, error) { return m.many, m.err }

type transformerMock struct{}
func (tr transformerMock) Execute(ctx context.Context, r *http.Request) (context.Context, *http.Request, error) { return ctx, r, nil }

type transformBuilderMock struct { many []pipeline.Transformer; err error }
func (m transformBuilderMock) BuildMany(config.Transform) ([]pipeline.Transformer, error) { return m.many, m.err }

type middlewareBuilderMock struct { many []pipeline.Middleware; err error }
func (m middlewareBuilderMock) Build(config.Middleware) (pipeline.Middleware, error) { return nil, nil }
func (m middlewareBuilderMock) BuildMany([]config.Middleware) ([]pipeline.Middleware, error) { return m.many, m.err }

type proxyMock struct{}
func (p proxyMock) ServeHTTP(http.ResponseWriter, *http.Request) {}

type proxyBuilderMock struct { got config.Upstream; err error }
func (m *proxyBuilderMock) Build(up config.Upstream) (pipeline.Proxy, error) { m.got = up; if m.err != nil { return nil, m.err }; return proxyMock{}, nil }

func TestEndpointBuilder_Build(t *testing.T) {
	pb := &proxyBuilderMock{}
	b := endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{}, pb)

	e, err := b.Build(config.Route{Path: "/v1", Method: "GET", Upstream: config.Upstream{URL: "http://example.com"}})
	if err != nil { t.Fatalf("expected nil error, got %v", err) }
	if e == nil { t.Fatal("expected endpoint") }
	if pb.got.Method != "GET" { t.Fatalf("expected upstream method fallback GET, got %s", pb.got.Method) }
}

func TestEndpointBuilder_Build_Errors(t *testing.T) {
	badErr := errors.New("bad")
	cases := []struct{ name string; builder *endpointBuilder.EndpointBuilder; cfg config.Route; exp error }{
		{"check error", endpointBuilder.NewBuilder(checkBuilderMock{err: badErr}, transformBuilderMock{}, middlewareBuilderMock{}, &proxyBuilderMock{}), config.Route{}, badErr},
		{"transform error", endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{err: badErr}, middlewareBuilderMock{}, &proxyBuilderMock{}), config.Route{}, badErr},
		{"middleware error", endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{err: badErr}, &proxyBuilderMock{}), config.Route{}, badErr},
		{"proxy error", endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{}, &proxyBuilderMock{err: badErr}), config.Route{Path: "/v1", Method: "GET", Upstream: config.Upstream{URL: "http://example.com", Method: "GET"}}, badErr},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := tc.builder.Build(tc.cfg)
			if !errors.Is(err, tc.exp) { t.Fatalf("expected %v got %v", tc.exp, err) }
			if e != nil { t.Fatal("expected nil endpoint") }
		})
	}
}

func TestEndpointBuilder_BuildMany(t *testing.T) {
	pb := &proxyBuilderMock{}
	b := endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{}, pb)

	res, err := b.BuildMany([]config.Route{{Path: "/v1", Method: "GET", Upstream: config.Upstream{URL: "http://example.com", Method: "GET"}}})
	if err != nil { t.Fatalf("expected nil error, got %v", err) }
	if len(res) != 1 { t.Fatalf("expected 1 endpoint, got %d", len(res)) }
}
