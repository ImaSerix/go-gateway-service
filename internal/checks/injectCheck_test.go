package checks_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestInjectCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.InjectCheckConfig
		expErr error
	}{
		{
			name: "success",
			cfg: &config.InjectCheckConfig{
				Ctx: map[string]any{
					"key":  "value",
					"key2": 21,
				},
			},
			expErr: nil,
		},
		{
			name:   "nil config",
			cfg:    nil,
			expErr: checks.ErrNilConfig,
		},
		{
			name: "empty inject context map",
			cfg: &config.InjectCheckConfig{
				Ctx: map[string]any{},
			},
			expErr: checks.ErrEmptyInjectContext,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewInjectCheck(test.cfg)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected wrapped error %v, but got %v", test.expErr, err)
			}
			if err == nil && c == nil {
				t.Fatal("got nil check, but no error")
			}
			if err != nil && c != nil {
				t.Fatal("got error, but check not nil")
			}
		})
	}
}

func TestInjectCheck_Execute_Success(t *testing.T) {
	cfg := &config.InjectCheckConfig{
		Ctx: map[string]any{
			"key":  "value",
			"key2": 21,
		},
	}

	c, err := checks.NewInjectCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)

	ctx, err := c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	v, ok := ctx.Value("key").(string)
	if !ok {
		t.Fatal("expected value be 'string' type")
	}
	if v != "value" {
		t.Fatalf("expected context value 'value', but got %s", v)
	}
	v2, ok := ctx.Value("key2").(int)
	if !ok {
		t.Fatal("expected value be 'string' type")
	}
	if v2 != 21 {
		t.Fatalf("expected context value 'value', but got %s", v)
	}

}

func TestInjectCheck_Execute_RequestNil(t *testing.T) {
	cfg := &config.InjectCheckConfig{
		Ctx: map[string]any{
			"key":  "value",
			"key2": 21,
		},
	}

	c, err := checks.NewInjectCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}
