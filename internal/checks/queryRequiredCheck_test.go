package checks_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestQueryRequiredCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.QueryRequiredCheckConfig
		expErr error
	}{
		{
			name: "success",
			cfg: &config.QueryRequiredCheckConfig{
				Queries: []string{
					"name",
					"category",
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
			name: "required headers list empty",
			cfg: &config.QueryRequiredCheckConfig{
				Queries: []string{},
			},
			expErr: checks.ErrEmptyQueries,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewQueryRequiredCheck(test.cfg)
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

func TestQueryRequiredCheck_Execute_Success(t *testing.T) {
	cfg := &config.QueryRequiredCheckConfig{
		Queries: []string{
			"name",
			"category",
		},
	}

	c, err := checks.NewQueryRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestQueryRequiredCheck_Execute_RequestNil(t *testing.T) {
	cfg := &config.QueryRequiredCheckConfig{
		Queries: []string{
			"name",
			"category",
		},
	}

	c, err := checks.NewQueryRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}

func TestQueryRequiredCheck_Execute_NoHeader(t *testing.T) {
	cfg := &config.QueryRequiredCheckConfig{
		Queries: []string{
			"name",
			"category",
		},
	}
	req := httptest.NewRequest("GET", "http://nice.url?name=name", nil)

	c, err := checks.NewQueryRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, checks.ErrNoQueryParam) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNoQueryParam, err)
	}
	if !strings.Contains(err.Error(), "category") {
		t.Fatalf("expected error contains 'X-Password', but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}
