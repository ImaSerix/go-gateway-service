package transformer

import "github.com/ImaSerix/go-gateway-service/internal/types"

type TransformRegistry struct {
	registry map[types.TransformerName]Factory
}

func NewTransformerRegistry() *TransformRegistry {
	return &TransformRegistry{
		registry: map[types.TransformerName]Factory{},
	}
}

func (cr *TransformRegistry) Get(key types.TransformerName) (Factory, bool) {
	f, ok := cr.registry[key]
	return f, ok
}

func (cr *TransformRegistry) Register(key types.TransformerName, factory Factory) {
	cr.registry[key] = factory
}
