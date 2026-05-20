package transformer

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
	"github.com/ImaSerix/go-gateway-service/internal/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/types"
)

type Factory interface {
	Create(raw any) (pipeline.Transformer, error)
}

type Registry interface {
	Get(key types.TransformerName) (Factory, bool)
}

type Builder struct {
	resolver resolver.Resolver
}

func NewBuilder(resolver resolver.Resolver) *Builder {
	return &Builder{
		resolver: resolver,
	}
}

func (b *Builder) BuildMany(cfg config.Transform) ([]pipeline.Transformer, error) {
	res := []pipeline.Transformer{}

	if cfg.Body != nil {
		res = append(res, transformer.NewBodyTransformer(cfg.Body, b.resolver))
	}

	if cfg.Header != nil {
		res = append(res, transformer.NewHeaderTransformer(cfg.Header, b.resolver))
	}

	return res, nil
}
