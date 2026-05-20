package builder

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
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
}

func RegisterChecks(registry *check.CheckRegistry, client *http.Client) {
	registry.Register(types.Auth, check.NewAuthFactory(client))
	registry.Register(types.HeaderRequired, check.NewHeaderRequiredFactory())
	registry.Register(types.Inject, check.NewInjectFactory())
	registry.Register(types.IPWhiteList, check.NewIPWhiteListFactory())
	registry.Register(types.QueryRequired, check.NewQueryRequiredFactory())
	registry.Register(types.RateLimitC, check.NewRateLimitFactory())
	registry.Register(types.TimeoutC, check.NewTimeoutFactory())
}

func RegisterResolvers(multiResolver *resolver.MultiResolver) {
	multiResolver.Register("ctx", resolver.NewContextResolver())
	multiResolver.Register("route", resolver.NewRouterResolver())
	multiResolver.Register("query", resolver.NewQueryResolver())
	multiResolver.Register("header", resolver.NewHeaderResolver())
}
