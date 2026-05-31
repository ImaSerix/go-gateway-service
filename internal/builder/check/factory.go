package check

import (
	"fmt"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"gopkg.in/yaml.v3"
)

// DEPRECATED
//
//	TODO: Нужно создать альтернативу - policy чек, с переиспользованием всяких нынешних хорошо работающих вещей
//
// type AuthFactory struct {
// 	client *http.Client
// }

// func NewAuthFactory(client *http.Client) *AuthFactory {
// 	return &AuthFactory{
// 		client: client,
// 	}
// }

// func (f *AuthFactory) Create(raw yaml.Node, render renderer.Renderer) (pipeline.Checker, error) {

// 	var cfg config.AuthCheck
// 	if err := raw.Decode(&cfg); err != nil {
// 		return nil, fmt.Errorf("new auth check factory: %w", err)
// 	}

// 	c, err := check.NewAuth(cfg, f.client)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c, nil
// }

type PolicyFactory struct {
	transformBuilder TransformBuilder
	clientBuilder    ClientBuilder
	storeBuilder     StoreBuilder
}

func NewPolicyFactory(t TransformBuilder, c ClientBuilder, s StoreBuilder) *PolicyFactory {
	return &PolicyFactory{
		transformBuilder: t,
		clientBuilder:    c,
		storeBuilder:     s,
	}
}

func (f *PolicyFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.PolicyCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new header required check factory: %w", err)
	}

	t, err := f.transformBuilder.BuildMany(cfg.Transform)
	if err != nil {
		return nil, fmt.Errorf("create policy check: %w", err)
	}

	c, err := f.clientBuilder.Build(cfg.Upstream)
	if err != nil {
		return nil, fmt.Errorf("create policy check: %w", err)
	}

	s, err := f.storeBuilder.Build(cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("create policy check: %w", err)
	}

	return check.NewPolicyCheck(t, c, s), nil
}

type HeaderRequiredFactory struct {
	render renderer.Renderer
}

func NewHeaderRequiredFactory(render renderer.Renderer) *HeaderRequiredFactory {
	return &HeaderRequiredFactory{
		render: render,
	}
}

func (f *HeaderRequiredFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.HeaderRequiredCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new header required check factory: %w", err)
	}

	c, err := check.NewHeaderRequiredCheck(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// DEPRECATED
//
// TODO: Создать альтернативный middleware, так как это не чек, он ничего не проверяет
//
// type InjectFactory struct{}

// func NewInjectFactory() *InjectFactory {
// 	return &InjectFactory{}
// }

// func (f *InjectFactory) Create(raw yaml.Node, render renderer.Renderer) (pipeline.Checker, error) {

// 	var cfg config.InjectCheck
// 	if err := raw.Decode(&cfg); err != nil {
// 		return nil, fmt.Errorf("new inject check factory: %w", err)
// 	}

// 	c, err := check.NewInject(cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c, nil
// }

type IPWhiteListFactory struct {
	render renderer.Renderer
}

func NewIPWhiteListFactory(render renderer.Renderer) *IPWhiteListFactory {
	return &IPWhiteListFactory{
		render: render,
	}
}

func (f *IPWhiteListFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.IPWhiteListCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new inject check factory: %w", err)
	}

	c, err := check.NewIPWhiteList(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type QueryRequiredFactory struct {
	render renderer.Renderer
}

func NewQueryRequiredFactory(render renderer.Renderer) *QueryRequiredFactory {
	return &QueryRequiredFactory{
		render: render,
	}
}

func (f *QueryRequiredFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.QueryRequiredCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new inject check factory: %w", err)
	}

	c, err := check.NewQueryRequired(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}
