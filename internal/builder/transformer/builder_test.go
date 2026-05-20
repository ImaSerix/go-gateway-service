package transformer_test

import (
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

func TestTransformerBuilder_BuildMany(t *testing.T) {
	b := transformer.NewBuilder(resolver.NewMultiResolver())

	res, err := b.BuildMany(config.Transform{Header: map[string]string{"X-Test": "header.req"}, Body: map[string]any{"user": "query.id"}})
	if err != nil { t.Fatalf("expected nil error, got %v", err) }
	if len(res) != 2 { t.Fatalf("expected 2 transformers, got %d", len(res)) }
}

func TestTransformerBuilder_BuildMany_Empty(t *testing.T) {
	b := transformer.NewBuilder(resolver.NewMultiResolver())

	res, err := b.BuildMany(config.Transform{})
	if err != nil { t.Fatalf("expected nil error, got %v", err) }
	if len(res) != 0 { t.Fatalf("expected 0 transformers, got %d", len(res)) }
}
