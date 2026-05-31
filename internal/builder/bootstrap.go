package builder

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/types"
)

func RegisterMiddlewares(registy *middleware.MiddlewareRegistry) {
	registy.Register(types.Cors, middleware.NewCorsFactory())
	registy.Register(types.Logging, middleware.NewLoggingFactory())
	registy.Register(types.Metric, middleware.NewMetricFactory())
	registy.Register(types.RateLimit, middleware.NewRateLimitFactory())
	registy.Register(types.RealIP, middleware.NewRealIPFactory())
	registy.Register(types.Recovery, middleware.NewRecoveryFactory())
	registy.Register(types.RequestID, middleware.NewRequestIDFactory())
	registy.Register(types.Timeout, middleware.NewTimeoutFactory())
	registy.Register(types.Inject, middleware.NewInjectFactory())
}

func RegisterChecks(registry *check.CheckRegistry, render renderer.Renderer, client *http.Client, t check.TransformBuilder, c check.ClientBuilder, s check.StoreBuilder) {
	// registry.Register(types.Auth, check.NewAuthFactory(client))
	registry.Register(types.Policy, check.NewPolicyFactory(t, c, s))
	registry.Register(types.HeaderRequired, check.NewHeaderRequiredFactory(render))
	registry.Register(types.IPWhiteList, check.NewIPWhiteListFactory(render))
	registry.Register(types.QueryRequired, check.NewQueryRequiredFactory(render))
}

func RegisterTransformers(registry *transformer.TransformRegistry, render renderer.Renderer) {
	registry.Register(types.Headers, transformer.NewHeadersFactory(render))
	registry.Register(types.BodyFields, transformer.NewBodyFieldsFactory(render))
	registry.Register(types.QueryParams, transformer.NewQueryParamsFactory(render))
}
