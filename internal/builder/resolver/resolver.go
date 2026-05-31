package resolver

import "github.com/ImaSerix/go-gateway-service/internal/resolver"

type MultiResolverBuilder struct{}

func NewMultiResolverBuilder() *MultiResolverBuilder {
	return &MultiResolverBuilder{}
}

func (b *MultiResolverBuilder) Build() *resolver.MultiResolver {
	return resolver.NewMultiResolver()
}
