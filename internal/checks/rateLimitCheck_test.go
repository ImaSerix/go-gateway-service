package checks_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestRateLimitCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.RateLimitCheckConfig
		expErr error
	}{
		{
			name: "success",
			cfg: &config.RateLimitCheckConfig{
				Limit:  50,
				Window: "1m",
			},
			expErr: nil,
		},
		{
			name:   "nil config",
			cfg:    nil,
			expErr: checks.ErrNilConfig,
		},
		{
			name: "invalid duration",
			cfg: &config.RateLimitCheckConfig{
				Limit:  50,
				Window: "bad duration",
			},
			expErr: checks.ErrInvalidWindow,
		},
		{
			name: "negative limit",
			cfg: &config.RateLimitCheckConfig{
				Limit:  -1,
				Window: "1m",
			},
			expErr: checks.ErrInvalidLimit,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewRateLimitCheck(test.cfg)
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

func TestRateLimitCheck_Execute_Success(t *testing.T) {
	cfg := &config.RateLimitCheckConfig{
		Limit:  50,
		Window: "1m",
	}

	c, err := checks.NewRateLimitCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
	req.RemoteAddr = "127.0.0.1"

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestRateLimitCheck_Execute_RequestNil(t *testing.T) {
	cfg := &config.RateLimitCheckConfig{
		Limit:  50,
		Window: "1m",
	}

	c, err := checks.NewRateLimitCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
}

func TestRateLimitCheck_Execute_TooManyRequests(t *testing.T) {
	cfg := &config.RateLimitCheckConfig{
		Limit:  5,
		Window: "1m",
	}
	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
	req.RemoteAddr = "127.0.0.1"

	c, err := checks.NewRateLimitCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	for i := 0; i < 5; i++ {
		ctx, err := c.Execute(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if err == nil && ctx == nil {
			t.Fatal("got nil check, but no error")
		}
		if err != nil && ctx != nil {
			t.Fatal("got error, but check not nil")
		}
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, checks.ErrTooManyRequests) {
		t.Fatalf("expected error %v, but got %v", checks.ErrTooManyRequests, err)
	}
}

func TestRateLimitCheck_Execute_LimitReset(t *testing.T) {
	cfg := &config.RateLimitCheckConfig{
		Limit:  5,
		Window: "1s",
	}
	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
	req.RemoteAddr = "127.0.0.1"

	c, err := checks.NewRateLimitCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	for i := 0; i < 5; i++ {
		_, err := c.Execute(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
	}

	time.Sleep(1100 * time.Millisecond)

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}
