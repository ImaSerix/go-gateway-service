package transformer

import (
	"fmt"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
	"github.com/ImaSerix/go-gateway-service/internal/transformer"
	"gopkg.in/yaml.v3"
)

type HeadersFactory struct {
	render renderer.Renderer
}

func NewHeadersFactory(render renderer.Renderer) *HeadersFactory {
	return &HeadersFactory{
		render: render,
	}
}

func (f *HeadersFactory) Create(raw yaml.Node) (pipeline.Transformer, error) {

	var cfg config.HeadersTransform
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new headers transformer factory: %w", err)
	}

	t := transformer.NewHeaderTransformer(cfg, f.render)

	return t, nil
}

type BodyFieldsFactory struct {
	render renderer.Renderer
}

func NewBodyFieldsFactory(render renderer.Renderer) *BodyFieldsFactory {
	return &BodyFieldsFactory{
		render: render,
	}
}

func (f *BodyFieldsFactory) Create(raw yaml.Node) (pipeline.Transformer, error) {

	var cfg config.BodyFieldsTransform
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new body fields transform factory: %w", err)
	}

	t := transformer.NewBodyTransformer(cfg, f.render)

	return t, nil
}

type QueryParamsFactory struct {
	render renderer.Renderer
}

func NewQueryParamsFactory(render renderer.Renderer) *QueryParamsFactory {
	return &QueryParamsFactory{
		render: render,
	}
}

func (f *QueryParamsFactory) Create(raw yaml.Node) (pipeline.Transformer, error) {

	var cfg config.QueryParamsTransform
	if err := raw.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("new query params transform factory: %w", err)
	}

	t := transformer.NewQueryParams(cfg, f.render)

	return t, nil
}
