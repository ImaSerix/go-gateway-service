package transformer_test

import (
	"errors"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder"
	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"gopkg.in/yaml.v3"
)

func TestTransformerBuilder_BuildMany(t *testing.T) {

	registry := transformer.NewTransformerRegistry()
	render := &rendererMock{}

	builder.RegisterTransformers(registry, render)

	b := transformer.NewBuilder(registry)

	res, err := b.BuildMany(config.Transform{
		"header":      yaml.Node{},
		"body_fields": yaml.Node{},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 transformers, got %d", len(res))
	}
}

func TestTransformerBuilder_BuildMany_UnregisteredTransformName(t *testing.T) {

	registry := transformer.NewTransformerRegistry()
	render := &rendererMock{}

	builder.RegisterTransformers(registry, render)

	b := transformer.NewBuilder(registry)

	res, err := b.BuildMany(config.Transform{
		"header":      yaml.Node{},
		"body_fields": yaml.Node{},
		"unknown":     yaml.Node{},
	})
	if !errors.Is(err, transformer.ErrUnregisteredTransformName) {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res != nil {
		t.Fatal("expected nil res, but got non-nil")
	}
}

func TestTransformerBuilder_BuildMany_Empty(t *testing.T) {
	registry := transformer.NewTransformerRegistry()
	render := &rendererMock{}

	builder.RegisterTransformers(registry, render)

	b := transformer.NewBuilder(registry)

	res, err := b.BuildMany(config.Transform{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 transformers, got %d", len(res))
	}
}
