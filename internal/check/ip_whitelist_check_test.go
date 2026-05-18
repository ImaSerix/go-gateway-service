package check_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestIPWhitelistCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.IPWhiteListCheck
		expErr error
	}{
		{
			name: "success",
			cfg: config.IPWhiteListCheck{
				IP: []string{
					"127.0.0.1",
					"192.168.0.1",
				},
			},
			expErr: nil,
		},
		{
			name: "ip whitelist empty",
			cfg: config.IPWhiteListCheck{
				IP: []string{},
			},
			expErr: check.ErrEmptyIP,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := check.NewIPWhiteList(test.cfg)
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

func TestIPWhitelistCheck_Execute_Success(t *testing.T) {
	cfg := config.IPWhiteListCheck{
		IP: []string{
			"127.0.0.1",
			"192.168.0.1",
		},
	}

	c, err := check.NewIPWhiteList(cfg)
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

func TestIPWhitelistCheck_Execute_WithPort(t *testing.T) {
	cfg := config.IPWhiteListCheck{
		IP: []string{
			"127.0.0.1",
			"192.168.0.1",
		},
	}

	c, err := check.NewIPWhiteList(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
	req.RemoteAddr = "127.0.0.1:55555"

	_, err = c.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestIPWhitelistCheck_Execute_RequestNil(t *testing.T) {
	cfg := config.IPWhiteListCheck{
		IP: []string{
			"127.0.0.1",
			"192.168.0.1",
		},
	}

	c, err := check.NewIPWhiteList(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	_, err = c.Execute(context.Background(), nil)
	if !errors.Is(err, check.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", check.ErrNilRequest, err)
	}
}

func TestIPWhitelistCheck_Execute_IPNotInList(t *testing.T) {
	cfg := config.IPWhiteListCheck{
		IP: []string{
			"127.0.0.1",
			"192.168.0.1",
		},
	}
	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
	req.RemoteAddr = ""

	c, err := check.NewIPWhiteList(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	_, err = c.Execute(context.Background(), req)
	if !errors.Is(err, check.ErrForbidden) {
		t.Fatalf("expected error %v, but got %v", check.ErrForbidden, err)
	}
}
