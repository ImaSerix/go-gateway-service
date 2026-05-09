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

func TestHeaderRequiredCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.HeaderRequiredCheckConfig
		expErr error
	}{
		{
			name: "success",
			cfg: &config.HeaderRequiredCheckConfig{
				Headers: []string{
					"X-Username",
					"X-Password",
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
			cfg: &config.HeaderRequiredCheckConfig{
				Headers: []string{},
			},
			expErr: checks.ErrEmptyHeaders,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewHeaderRequiredCheck(test.cfg)
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

func TestHeaderRequiredCheck_Execute_Success(t *testing.T) {
	cfg := &config.HeaderRequiredCheckConfig{
		Headers: []string{
			"X-Username",
			"X-Password",
		},
	}

	c, err := checks.NewHeaderRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url", nil)
	req.Header.Set("X-Username", "nice username")
	req.Header.Set("X-Password", "nice password")

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestHeaderRequiredCheck_Execute_RequestNil(t *testing.T) {
	cfg := &config.HeaderRequiredCheckConfig{
		Headers: []string{
			"X-Username",
			"X-Password",
		},
	}

	c, err := checks.NewHeaderRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}

func TestHeaderRequiredCheck_Execute_NoHeader(t *testing.T) {
	cfg := &config.HeaderRequiredCheckConfig{
		Headers: []string{
			"X-Username",
			"X-Password",
		},
	}

	req := httptest.NewRequest("GET", "http://nice.url", nil)
	req.Header.Set("X-Username", "nice username")

	c, err := checks.NewHeaderRequiredCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, checks.ErrNoHeader) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNoHeader, err)
	}
	if !strings.Contains(err.Error(), "X-Password") {
		t.Fatalf("expected error contains 'X-Password', but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}
