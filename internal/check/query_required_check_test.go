package check_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestQueryRequiredCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.QueryRequiredCheck
		expErr error
	}{
		{
			name: "success",
			cfg: config.QueryRequiredCheck{
				QueryParams: []string{
					"name",
					"category",
				},
			},
			expErr: nil,
		},
		{
			name: "required headers list empty",
			cfg: config.QueryRequiredCheck{
				QueryParams: []string{},
			},
			expErr: check.ErrEmptyQuery,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := check.NewQueryRequired(test.cfg)
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
	cfg := config.QueryRequiredCheck{
		QueryParams: []string{
			"name",
			"category",
		},
	}

	c, err := check.NewQueryRequired(cfg)
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
	cfg := config.QueryRequiredCheck{
		QueryParams: []string{
			"name",
			"category",
		},
	}

	c, err := check.NewQueryRequired(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, check.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", check.ErrNilRequest, err)
	}
}

func TestQueryRequiredCheck_Execute_NoHeader(t *testing.T) {
	cfg := config.QueryRequiredCheck{
		QueryParams: []string{
			"name",
			"category",
		},
	}
	req := httptest.NewRequest("GET", "http://nice.url?name=name", nil)

	c, err := check.NewQueryRequired(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, check.ErrNoQueryParam) {
		t.Fatalf("expected error %v, but got %v", check.ErrNoQueryParam, err)
	}
	if !strings.Contains(err.Error(), "category") {
		t.Fatalf("expected error contains 'X-Password', but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, check.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", check.ErrNilRequest, err)
	}
}
