package transformer

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

type Builder struct {
	registry Registry
}

func NewBuilder(registry Registry) *Builder {
	return &Builder{
		registry: registry,
	}
}

func (b *Builder) Build(t string, raw yaml.Node) (pipeline.Transformer, error) {

	f, ok := b.registry.Get(types.TransformerName(t))
	if !ok {
		return nil, ErrUnregisteredTransformName
	}

	tr, err := f.Create(raw)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

func (b *Builder) BuildMany(cfg config.Transform) ([]pipeline.Transformer, error) {

	res := []pipeline.Transformer{}

	for t, v := range cfg {
		tr, err := b.Build(t, v)
		if err != nil {
			return nil, err
		}
		res = append(res, tr)
	}

	return res, nil
}
