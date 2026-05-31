package store

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/store"
)

type Builder struct {
	render renderer.ResponseRenderer
}

func NewBuilder(render renderer.ResponseRenderer) *Builder {
	return &Builder{
		render: render,
	}
}

func (b *Builder) Build(cfg config.Store) (pipeline.Store, error) {
	return store.NewStore(cfg, b.render), nil
}
