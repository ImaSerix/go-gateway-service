package check

import (
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type AuthFactory struct {
	client *http.Client
}

func NewAuthFactory(client *http.Client) *AuthFactory {
	return &AuthFactory{
		client: client,
	}
}

func (f *AuthFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.AuthCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new auth check factory: %w", err)
	}

	c, err := check.NewAuth(cfg, f.client)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type HeaderRequiredFactory struct{}

func NewHeaderRequiredFactory() *HeaderRequiredFactory {
	return &HeaderRequiredFactory{}
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

type InjectFactory struct{}

func NewInjectFactory() *InjectFactory {
	return &InjectFactory{}
}

func (f *InjectFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.InjectCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new inject check factory: %w", err)
	}

	c, err := check.NewInject(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type IPWhiteListFactory struct{}

func NewIPWhiteListFactory() *IPWhiteListFactory {
	return &IPWhiteListFactory{}
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

type QueryRequiredFactory struct{}

func NewQueryRequiredFactory() *QueryRequiredFactory {
	return &QueryRequiredFactory{}
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

type RateLimitFactory struct{}

func NewRateLimitFactory() *RateLimitFactory {
	return &RateLimitFactory{}
}

func (f *RateLimitFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.RateLimitCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new inject check factory: %w", err)
	}

	c, err := check.NewRateLimit(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type TimeoutFactory struct{}

func NewTimeoutFactory() *TimeoutFactory {
	return &TimeoutFactory{}
}

func (f *TimeoutFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.TimeoutCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new inject check factory: %w", err)
	}

	c, err := check.NewTimeout(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}
