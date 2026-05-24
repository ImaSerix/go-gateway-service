package endpoint_test

import (
	"errors"
	"testing"

	endpointBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/endpoint"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestEndpointBuilder_Build(t *testing.T) {
	pb := &proxyBuilderMock{}
	b := endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{}, pb)

	e, err := b.Build(
		config.Route{
			Path:   "/v1",
			Method: "GET",
			Upstream: config.Upstream{
				Host:   "example.com",
				Scheme: "http",
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if e == nil {
		t.Fatal("expected endpoint")
	}
	if pb.got.Method != "GET" {
		t.Fatalf("expected upstream method fallback GET, got %s", pb.got.Method)
	}
}

func TestEndpointBuilder_Build_Errors(t *testing.T) {
	badErr := errors.New("bad")
	cases := []struct {
		name    string
		builder *endpointBuilder.EndpointBuilder
		cfg     config.Route
		exp     error
	}{
		{
			name: "check error",
			builder: endpointBuilder.NewBuilder(
				checkBuilderMock{err: badErr},
				transformBuilderMock{},
				middlewareBuilderMock{},
				&proxyBuilderMock{}),
			cfg: config.Route{},
			exp: badErr,
		},
		{
			name: "transform error",
			builder: endpointBuilder.NewBuilder(
				checkBuilderMock{},
				transformBuilderMock{err: badErr},
				middlewareBuilderMock{},
				&proxyBuilderMock{}),
			cfg: config.Route{},
			exp: badErr,
		},
		{
			name: "transform error",
			builder: endpointBuilder.NewBuilder(
				checkBuilderMock{},
				transformBuilderMock{err: badErr},
				middlewareBuilderMock{},
				&proxyBuilderMock{}),
			cfg: config.Route{},
			exp: badErr,
		},
		{
			name: "middleware error",
			builder: endpointBuilder.NewBuilder(
				checkBuilderMock{},
				transformBuilderMock{},
				middlewareBuilderMock{err: badErr},
				&proxyBuilderMock{}),
			cfg: config.Route{},
			exp: badErr,
		},
		{
			name: "proxy error",
			builder: endpointBuilder.NewBuilder(
				checkBuilderMock{},
				transformBuilderMock{},
				middlewareBuilderMock{},
				&proxyBuilderMock{err: badErr}),
			cfg: config.Route{
				Path:   "/v1",
				Method: "GET",
				Upstream: config.Upstream{
					Host:   "example.com",
					Scheme: "http",
					Method: "GET",
				},
			},
			exp: badErr,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := tc.builder.Build(tc.cfg)
			if !errors.Is(err, tc.exp) {
				t.Fatalf("expected %v got %v", tc.exp, err)
			}
			if e != nil {
				t.Fatal("expected nil endpoint")
			}
		})
	}
}

func TestEndpointBuilder_BuildMany(t *testing.T) {
	pb := &proxyBuilderMock{}
	b := endpointBuilder.NewBuilder(checkBuilderMock{}, transformBuilderMock{}, middlewareBuilderMock{}, pb)

	res, err := b.BuildMany([]config.Route{
		{
			Path:   "/v1",
			Method: "GET",
			Upstream: config.Upstream{
				Host:   "example.com",
				Scheme: "http",
				Method: "GET",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(res))
	}
}
