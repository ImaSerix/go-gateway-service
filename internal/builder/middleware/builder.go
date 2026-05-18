package middleware

import (
	"fmt"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

type Factory interface {
	Create(raw yaml.Node) (pipeline.Middleware, error)
}

type Registry interface {
	Get(key types.MiddlewareName) (Factory, bool)
}

type Builder struct {
	registry Registry
}

func NewBuilder(registry Registry) *Builder {
	return &Builder{
		registry: registry,
	}
}

func (b *Builder) Build(cfg config.Middleware) (pipeline.Middleware, error) {

	f, ok := b.registry.Get(types.MiddlewareName(cfg.Type))
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnregisteredMiddlewareType, cfg.Type)
	}

	m, err := f.Create(cfg.Config)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (b *Builder) BuildMany(cfg []config.Middleware) ([]pipeline.Middleware, error) {
	result := make([]pipeline.Middleware, 0, len(cfg))

	for _, c := range cfg {
		m, err := b.Build(c)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}
