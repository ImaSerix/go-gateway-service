package middleware_test

import (
	"testing"

	middleawreBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"gopkg.in/yaml.v3"
)

func TestCheckFactory(t *testing.T) {
	tests := []struct {
		name    string
		factory middleawreBuilder.Factory
		cfg     any
		wantErr bool
	}{
		{
			name:    "cors",
			factory: middleawreBuilder.NewCorsFactory(),
			cfg: config.CORSMiddleaware{
				Allowed: config.CORSMiddleaware_Allowed{
					Origin: []string{
						"*",
					},
				},
			},
		},
		{
			name:    "cors error",
			factory: middleawreBuilder.NewCorsFactory(),
			cfg:     config.CORSMiddleaware{},
			wantErr: true,
		},
		{
			name:    "logging",
			factory: middleawreBuilder.NewLoggingFactory(),
			cfg:     nil,
		},
		{
			name:    "metric",
			factory: middleawreBuilder.NewMetricFactory(),
			cfg:     nil,
		},
		{
			name:    "rate limit",
			factory: middleawreBuilder.NewRateLimitFactory(),
			cfg: config.RateLimitMiddleware{
				Limit:  50,
				Window: "2s",
			},
		},
		{
			name:    "rate limit error",
			factory: middleawreBuilder.NewRateLimitFactory(),
			cfg:     config.RateLimitMiddleware{},
			wantErr: true,
		},
		{
			name:    "real ip",
			factory: middleawreBuilder.NewRealIPFactory(),
			cfg:     nil,
		},

		{
			name:    "recovery",
			factory: middleawreBuilder.NewRecoveryFactory(),
			cfg:     nil,
		},
		{
			name:    "request id",
			factory: middleawreBuilder.NewRequestIDFactory(),
			cfg:     nil,
		},

		{
			name:    "timeout",
			factory: middleawreBuilder.NewTimeoutFactory(),
			cfg: config.TimeoutMiddleware{
				Duration: "2s",
			},
		},
		{
			name:    "timeout error",
			factory: middleawreBuilder.NewTimeoutFactory(),
			cfg:     config.TimeoutMiddleware{},
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

			if got == nil {
				t.Fatal("expected non-nil midldeware, but got nil")
			}

		})
	}
}
