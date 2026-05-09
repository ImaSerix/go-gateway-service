package checks_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestTimeoutCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.TimeoutCheckConfig
		expErr error
	}{
		{
			name: "success",
			cfg: &config.TimeoutCheckConfig{
				Duration: "2s",
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
			cfg: &config.TimeoutCheckConfig{
				Duration: "invalid duration",
			},
			expErr: checks.ErrInvalidDuration,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewTimeoutCheck(test.cfg)
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

func TestTimeoutCheck_Execute_Success(t *testing.T) {
	cfg := &config.TimeoutCheckConfig{
		Duration: "2s",
	}

	c, err := checks.NewTimeoutCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestTimeoutCheck_Execute_DeadlineExceeded(t *testing.T) {
	cfg := &config.TimeoutCheckConfig{
		Duration: "200ms",
	}

	c, err := checks.NewTimeoutCheck(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)

	ctx, err := c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("expected context done with error %v, but got %v", ctx.Err(), context.DeadlineExceeded)
	}
}
