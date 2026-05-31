package check_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	checkBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/check"
	check "github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/client"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type transformBuilderMock struct{}

func (m transformBuilderMock) BuildMany(config.Transform) ([]pipeline.Transformer, error) {
	return nil, nil
}

type clientBuilderMock struct{}

func (m clientBuilderMock) Build(config.Upstream) (*client.Upstream, error) {
	return client.NewUpstreamClient(http.DefaultClient, &rendererMock{}, nil, http.MethodGet), nil
}

type storeBuilderMock struct{}

func (m storeBuilderMock) Build(config.Store) (pipeline.Store, error) {
	return storeMock{}, nil
}

type storeMock struct{}

func (m storeMock) Save(ctx context.Context, r *http.Response) (context.Context, error) {
	return ctx, nil
}

func TestCheckFactory(t *testing.T) {
	tests := []struct {
		name     string
		factory  checkBuilder.Factory
		cfg      any
		wantType any
		wantErr  bool
	}{
		{
			name:    "policy",
			factory: checkBuilder.NewPolicyFactory(transformBuilderMock{}, clientBuilderMock{}, storeBuilderMock{}),
			cfg: config.PolicyCheck{
				Upstream: config.Upstream{
					Scheme: "http",
					Host:   "policy.local",
					Method: http.MethodGet,
				},
				ExpectedStatus: http.StatusOK,
			},
			wantType: &check.PolicyCheck{},
		},
		{
			name:    "header required",
			factory: checkBuilder.NewHeaderRequiredFactory(),
			cfg: config.HeaderRequiredCheck{
				Headers: []string{
					"X-Username",
				},
			},
			wantType: &check.HeaderRequired{},
		},
		{
			name:    "header required error",
			factory: checkBuilder.NewHeaderRequiredFactory(),
			cfg:     config.HeaderRequiredCheck{},
			wantErr: true,
		},

		{
			name:    "ip whitelist",
			factory: checkBuilder.NewIPWhiteListFactory(),
			cfg: config.IPWhiteListCheck{
				IPs: []string{
					"127.0.0.1",
				},
			},
			wantType: &check.IPWhiteList{},
		},
		{
			name:    "ip whitelist error",
			factory: checkBuilder.NewIPWhiteListFactory(),
			cfg:     config.IPWhiteListCheck{},
			wantErr: true,
		},

		{
			name:    "query required",
			factory: checkBuilder.NewQueryRequiredFactory(),
			cfg: config.QueryRequiredCheck{
				QueryParams: []string{
					"limit",
				},
			},
			wantType: &check.QueryRequired{},
		},
		{
			name:    "query required error",
			factory: checkBuilder.NewQueryRequiredFactory(),
			cfg:     config.QueryRequiredCheck{},
			wantErr: true,
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
				t.Fatalf("expected check with type %T, but got %t", test.wantType, got)
			}

		})
	}
}
