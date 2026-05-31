package check

import (
	"fmt"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"gopkg.in/yaml.v3"
)

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
		return nil, fmt.Errorf("new policy check factory: %w", err)
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

	return check.NewPolicyCheck(t, c, s, cfg.ExpectedStatus), nil
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

type IPWhiteListFactory struct{}

func NewIPWhiteListFactory() *IPWhiteListFactory {
	return &IPWhiteListFactory{}
}

func (f *IPWhiteListFactory) Create(raw yaml.Node) (pipeline.Checker, error) {

	var cfg config.IPWhiteListCheck
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new ip whitelist check factory: %w", err)
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
		return nil, fmt.Errorf("new query required check factory: %w", err)
	}

	c, err := check.NewQueryRequired(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}
