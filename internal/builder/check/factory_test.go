package check_test

import (
	"net/http"
	"reflect"
	"testing"

	checkBuilder "github.com/ImaSerix/go-gateway-service/internal/builder/check"
	check "github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"gopkg.in/yaml.v3"
)

func TestCheckFactory(t *testing.T) {
	tests := []struct {
		name     string
		factory  checkBuilder.Factory
		cfg      any
		wantType any
		wantErr  bool
	}{
		{
			name:    "auth",
			factory: checkBuilder.NewAuthFactory(http.DefaultClient),
			cfg: config.AuthCheck{
				URL:    "http://nice.url",
				Method: "GET",
			},
			wantType: &check.Auth{},
		},
		{
			name:    "auth error",
			factory: checkBuilder.NewAuthFactory(http.DefaultClient),
			cfg:     config.AuthCheck{},
			wantErr: true,
		},

		{
			name:    "header required",
			factory: checkBuilder.NewHeaderRequiredFactory(),
			cfg: config.HeaderRequiredCheck{
				Header: []string{
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
			name:    "inject",
			factory: checkBuilder.NewInjectFactory(),
			cfg: config.InjectCheck{
				Ctx: map[string]any{
					"key": "something",
				},
			},
			wantType: &check.Inject{},
		},
		{
			name:    "inject error",
			factory: checkBuilder.NewInjectFactory(),
			cfg:     config.InjectCheck{},
			wantErr: true,
		},

		{
			name:    "ip whitelist",
			factory: checkBuilder.NewIPWhiteListFactory(),
			cfg: config.IPWhiteListCheck{
				IP: []string{
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
				Query: []string{
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

		{
			name:    "rate limit",
			factory: checkBuilder.NewRateLimitFactory(),
			cfg: config.RateLimitCheck{
				Limit:  50,
				Window: "1m",
			},
			wantType: &check.RateLimit{},
		},
		{
			name:    "rate limit error",
			factory: checkBuilder.NewRateLimitFactory(),
			cfg:     config.RateLimitCheck{},
			wantErr: true,
		},

		{
			name:    "timeout",
			factory: checkBuilder.NewTimeoutFactory(),
			cfg: config.TimeoutCheck{
				Duration: "2s",
			},
			wantType: &check.Timeout{},
		},
		{
			name:    "rate limit error",
			factory: checkBuilder.NewTimeoutFactory(),
			cfg:     config.TimeoutCheck{},
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
