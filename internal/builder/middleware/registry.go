package middleware

import "github.com/ImaSerix/go-gateway-service/internal/types"

type MiddlewareRegistry struct {
	registry map[types.MiddlewareName]Factory
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		registry: map[types.MiddlewareName]Factory{},
	}
}

func (mr *MiddlewareRegistry) Get(key types.MiddlewareName) (Factory, bool) {
	f, ok := mr.registry[key]
	return f, ok
}

func (mr *MiddlewareRegistry) Register(key types.MiddlewareName, factory Factory) {
	mr.registry[key] = factory
}
