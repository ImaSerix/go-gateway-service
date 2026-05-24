package middleware

import (
	"fmt"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type CorsFactory struct{}

func NewCorsFactory() *CorsFactory {
	return &CorsFactory{}
}

func (f *CorsFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {

	var cfg config.CORSMiddleaware

	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new cors middleware factory: %w", err)
	}

	if len(cfg.Allowed.Origin) == 0 {
		return nil, ErrEmptyAllowedOrigins
	}

	return middleware.CORS(cfg.Allowed.Origin, cfg.Allowed.Method, cfg.Allowed.Header), nil
}

type LoggingFactory struct{}

func NewLoggingFactory() *LoggingFactory {
	return &LoggingFactory{}
}

func (f *LoggingFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {
	return middleware.Logging, nil
}

type MetricFactory struct{}

func NewMetricFactory() *MetricFactory {
	return &MetricFactory{}
}

func (f *MetricFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {
	m := middleware.NewMetric()
	return m.Middleware(), nil
}

type RateLimitFactory struct{}

func NewRateLimitFactory() *RateLimitFactory {
	return &RateLimitFactory{}
}

func (f *RateLimitFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {

	var cfg config.RateLimitMiddleware

	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new rate limit middleware factory: %w", err)
	}

	if cfg.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	w, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, ErrInvalidWindow
	}

	m := middleware.NewRateLimit(cfg.Limit, w)
	return m.Middleware(), nil
}

type RealIPFactory struct{}

func NewRealIPFactory() *RealIPFactory {
	return &RealIPFactory{}
}

func (f *RealIPFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {
	return middleware.RealIP, nil
}

type RecoveryFactory struct{}

func NewRecoveryFactory() *RecoveryFactory {
	return &RecoveryFactory{}
}

func (f *RecoveryFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {
	return middleware.Recovery, nil
}

type RequestIDFactory struct{}

func NewRequestIDFactory() *RequestIDFactory {
	return &RequestIDFactory{}
}

func (f *RequestIDFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {
	return middleware.RequestID, nil
}

type TimeoutFactory struct{}

func NewTimeoutFactory() *TimeoutFactory {
	return &TimeoutFactory{}
}

func (f *TimeoutFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {

	var cfg config.TimeoutMiddleware

	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new cors middleware factory: %w", err)
	}

	d, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, ErrInvalidDuration
	}

	return middleware.Timeout(d), nil
}

type InjectFactory struct{}

func NewInjectFactory() *InjectFactory {
	return &InjectFactory{}
}

func (f *InjectFactory) Create(raw yaml.Node) (pipeline.Middleware, error) {

	var cfg config.InjectMiddleware

	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new cors middleware factory: %w", err)
	}

	return middleware.Inject(cfg.Context), nil
}
