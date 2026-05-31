package store

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Builder struct {
	render renderer.Renderer
}

func NewBuilder(render renderer.Renderer) *Builder {
	return &Builder{
		render: render,
	}
}

func (b *Builder) Build(cfg config.Store)
