package resolver

import "github.com/ImaSerix/go-gateway-service/internal/resolver"

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *resolver.MultiResolver {
	return resolver.NewMultiResolver()
}
