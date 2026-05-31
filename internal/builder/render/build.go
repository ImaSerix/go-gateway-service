package render

import (
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Builder struct {
	resolver renderer.Resolver
}

func NewBuilder(resolver renderer.Resolver) *Builder {
	return &Builder{
		resolver: resolver,
	}
}

func (b *Builder) Build() renderer.Renderer {
	return renderer.NewRender(b.resolver)
}
