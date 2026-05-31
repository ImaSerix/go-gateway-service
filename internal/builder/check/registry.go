package check

import "github.com/ImaSerix/go-gateway-service/internal/types"

type CheckRegistry struct {
	registry map[types.CheckName]Factory
}

func NewCheckRegistry() *CheckRegistry {
	return &CheckRegistry{
		registry: map[types.CheckName]Factory{},
	}
}

func (cr *CheckRegistry) Get(key types.CheckName) (Factory, bool) {
	f, ok := cr.registry[key]
	return f, ok
}

func (cr *CheckRegistry) Register(key types.CheckName, factory Factory) {
	cr.registry[key] = factory
}
