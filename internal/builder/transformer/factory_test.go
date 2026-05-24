package transformer_test

import (
	"reflect"
	"testing"

	transformerBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/transformer"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/transformer"
	"gopkg.in/yaml.v3"
)

func TestTransformerFactory(t *testing.T) {

	render := &rendererMock{}

	tests := []struct {
		name     string
		factory  transformerBuilder.Factory
		cfg      any
		wantType any
		wantErr  bool
	}{
		{
			name:     "headers",
			factory:  transformerBuilder.NewHeadersFactory(render),
			cfg:      config.HeadersTransform{},
			wantType: &transformer.Headers{},
		},
		{
			name:     "body_fields",
			factory:  transformerBuilder.NewBodyFieldsFactory(render),
			cfg:      config.BodyFieldsTransform{},
			wantType: &transformer.BodyFields{},
		},
		{
			name:     "query_params",
			factory:  transformerBuilder.NewQueryParamsFactory(render),
			cfg:      config.QueryParamsTransform{},
			wantType: &transformer.QueryParams{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var raw yaml.Node

			if err := raw.Encode(test.cfg); err != nil {
				t.Fatal(err)
			}

			got, err := test.factory.Create(raw)

			if test.wantErr {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, but got %v", err)
			}

			if reflect.TypeOf(got) != reflect.TypeOf(test.wantType) {
				t.Fatalf("expected transformer with type %T, but got %t", test.wantType, got)
			}
		})
	}
}
