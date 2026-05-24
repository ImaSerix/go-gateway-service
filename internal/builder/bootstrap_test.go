package builder_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/types"
)

type rendererMock struct {
	res string
	err error
}

func (rm *rendererMock) Render(s string, r *http.Request) (string, error) {
	return rm.res, rm.err
}

func TestBootstrap_RegisterMiddlewares(t *testing.T) {
	tests := []struct {
		name    string
		key     types.MiddlewareName
		expType any
	}{
		{
			name:    "cors",
			key:     types.Cors,
			expType: &middleware.CorsFactory{},
		},
		{
			name:    "logging",
			key:     types.Logging,
			expType: &middleware.LoggingFactory{},
		},
		{
			name:    "metric",
			key:     types.Metric,
			expType: &middleware.MetricFactory{},
		},
		{
			name:    "rate limit",
			key:     types.RateLimit,
			expType: &middleware.RateLimitFactory{},
		},
		{
			name:    "real ip",
			key:     types.RealIP,
			expType: &middleware.RealIPFactory{},
		},
		{
			name:    "recovery",
			key:     types.Recovery,
			expType: &middleware.RecoveryFactory{},
		},
		{
			name:    "request",
			key:     types.RequestID,
			expType: &middleware.RequestIDFactory{},
		},
		{
			name:    "timeout",
			key:     types.Timeout,
			expType: &middleware.TimeoutFactory{},
		},
		{
			name:    "inject",
			key:     types.Inject,
			expType: &middleware.InjectFactory{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := middleware.NewMiddlewareRegistry()
			builder.RegisterMiddlewares(reg)

			f, ok := reg.Get(test.key)
			if !ok {
				t.Fatalf("expected registry having this key, but it is not")
			}

			if reflect.TypeOf(test.expType) != reflect.TypeOf(f) {
				t.Fatalf("expected type %T, but got %T", test.expType, f)
			}
		})
	}
}

func TestBootstrap_RegisterChecks(t *testing.T) {
	tests := []struct {
		name    string
		key     types.CheckName
		expType any
	}{
		// {
		// 	name:    "aut",
		// 	key:     types.Auth,
		// 	expType: &check.AuthFactory{},
		// },
		{
			name:    "header required",
			key:     types.HeaderRequired,
			expType: &check.HeaderRequiredFactory{},
		},
		// {
		// 	name:    "inject",
		// 	key:     types.Inject,
		// 	expType: &check.InjectFactory{},
		// },
		{
			name:    "ip whitelist",
			key:     types.IPWhiteList,
			expType: &check.IPWhiteListFactory{},
		},
		{
			name:    "query required",
			key:     types.QueryRequired,
			expType: &check.QueryRequiredFactory{},
		},
		// {
		// 	name:    "rate limit",
		// 	key:     types.RateLimitC,
		// 	expType: &check.RateLimitFactory{},
		// },
		// {
		// 	name:    "timeout",
		// 	key:     types.TimeoutC,
		// 	expType: &check.TimeoutFactory{},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			render := &rendererMock{}

			reg := check.NewCheckRegistry()
			builder.RegisterChecks(reg, render, http.DefaultClient)

			f, ok := reg.Get(test.key)
			if !ok {
				t.Fatalf("expected registry having this key, but it is not")
			}

			if reflect.TypeOf(test.expType) != reflect.TypeOf(f) {
				t.Fatalf("expected type %T, but got %T", test.expType, f)
			}
		})
	}
}

func TestBootstrap_RegisterTransformer(t *testing.T) {
	tests := []struct {
		name    string
		key     types.TransformerName
		expType any
	}{
		{
			name:    "headers",
			key:     types.Headers,
			expType: &transformer.HeadersFactory{},
		},
		{
			name:    "body_fields",
			key:     types.BodyFields,
			expType: &transformer.BodyFieldsFactory{},
		},
		{
			name:    "query_params",
			key:     types.QueryParams,
			expType: &transformer.QueryParamsFactory{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			render := &rendererMock{}

			reg := transformer.NewTransformerRegistry()
			builder.RegisterTransformers(reg, render)

			f, ok := reg.Get(test.key)
			if !ok {
				t.Fatalf("expected registry having this key, but it is not")
			}

			if reflect.TypeOf(test.expType) != reflect.TypeOf(f) {
				t.Fatalf("expected type %T, but got %T", test.expType, f)
			}
		})
	}
}
