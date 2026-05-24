package endpoint_test

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type checkerMock struct{}

func (c checkerMock) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}

type checkBuilderMock struct {
	many []pipeline.Checker
	err  error
}

func (m checkBuilderMock) Build(config.Check) (pipeline.Checker, error) {
	return checkerMock{}, nil
}
func (m checkBuilderMock) BuildMany([]config.Check) ([]pipeline.Checker, error) {
	return m.many, m.err
}

type transformerMock struct{}

func (tr transformerMock) Execute(ctx context.Context, r *http.Request) (context.Context, *http.Request, error) {
	return ctx, r, nil
}

type transformBuilderMock struct {
	many []pipeline.Transformer
	err  error
}

func (m transformBuilderMock) BuildMany(config.Transform) ([]pipeline.Transformer, error) {
	return m.many, m.err
}

type middlewareBuilderMock struct {
	many []pipeline.Middleware
	err  error
}

func (m middlewareBuilderMock) Build(config.Middleware) (pipeline.Middleware, error) {
	return nil, nil
}
func (m middlewareBuilderMock) BuildMany([]config.Middleware) ([]pipeline.Middleware, error) {
	return m.many, m.err
}

type proxyMock struct{}

func (p proxyMock) ServeHTTP(http.ResponseWriter, *http.Request) {}

type proxyBuilderMock struct {
	got config.Upstream
	err error
}

func (m *proxyBuilderMock) Build(up config.Upstream) (pipeline.Proxy, error) {
	m.got = up
	if m.err != nil {
		return nil, m.err
	}
	return proxyMock{}, nil
}
