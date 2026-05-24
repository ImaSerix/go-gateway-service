package check_test

// import (
// 	"context"
// 	"errors"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/ImaSerix/go-gateway-service/internal/check"
// 	"github.com/ImaSerix/go-gateway-service/internal/config"
// )

// func TestRateLimitCheck(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		cfg    config.RateLimitCheck
// 		expErr error
// 	}{
// 		{
// 			name: "success",
// 			cfg: config.RateLimitCheck{
// 				Limit:  50,
// 				Window: "1m",
// 			},
// 			expErr: nil,
// 		},
// 		{
// 			name: "invalid duration",
// 			cfg: config.RateLimitCheck{
// 				Limit:  50,
// 				Window: "bad duration",
// 			},
// 			expErr: check.ErrInvalidWindow,
// 		},
// 		{
// 			name: "invalid duration: 0 duration",
// 			cfg: config.RateLimitCheck{
// 				Limit:  50,
// 				Window: "0s",
// 			},
// 			expErr: check.ErrInvalidWindow,
// 		},
// 		{
// 			name: "negative limit",
// 			cfg: config.RateLimitCheck{
// 				Limit:  -1,
// 				Window: "1m",
// 			},
// 			expErr: check.ErrInvalidLimit,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			c, err := check.NewRateLimit(test.cfg)
// 			if !errors.Is(err, test.expErr) {
// 				t.Fatalf("expected wrapped error %v, but got %v", test.expErr, err)
// 			}
// 			if err == nil && c == nil {
// 				t.Fatal("got nil check, but no error")
// 			}
// 			if err != nil && c != nil {
// 				t.Fatal("got error, but check not nil")
// 			}
// 		})
// 	}
// }

// func TestRateLimitCheck_Execute_Success(t *testing.T) {
// 	cfg := config.RateLimitCheck{
// 		Limit:  50,
// 		Window: "1m",
// 	}

// 	c, err := check.NewRateLimit(cfg)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}

// 	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
// 	req.RemoteAddr = "127.0.0.1"

// 	_, err = c.Execute(context.Background(), req)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}
// }

// func TestRateLimitCheck_Execute_RequestNil(t *testing.T) {
// 	cfg := config.RateLimitCheck{
// 		Limit:  50,
// 		Window: "1m",
// 	}

// 	c, err := check.NewRateLimit(cfg)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}
// 	_, err = c.Execute(context.Background(), nil)
// 	if !errors.Is(err, check.ErrNilRequest) {
// 		t.Fatalf("expected error %v, but got %v", check.ErrNilRequest, err)
// 	}
// }

// func TestRateLimitCheck_Execute_TooManyRequests(t *testing.T) {
// 	cfg := config.RateLimitCheck{
// 		Limit:  5,
// 		Window: "1m",
// 	}
// 	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
// 	req.RemoteAddr = "127.0.0.1"

// 	c, err := check.NewRateLimit(cfg)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}

// 	for i := 0; i < 5; i++ {
// 		ctx, err := c.Execute(context.Background(), req)
// 		if err != nil {
// 			t.Fatalf("expected no error, but got %v", err)
// 		}
// 		if err == nil && ctx == nil {
// 			t.Fatal("got nil check, but no error")
// 		}
// 		if err != nil && ctx != nil {
// 			t.Fatal("got error, but check not nil")
// 		}
// 	}

// 	_, err = c.Execute(context.Background(), req)
// 	if !errors.Is(err, check.ErrTooManyRequests) {
// 		t.Fatalf("expected error %v, but got %v", check.ErrTooManyRequests, err)
// 	}
// }

// func TestRateLimitCheck_Execute_LimitReset(t *testing.T) {
// 	cfg := config.RateLimitCheck{
// 		Limit:  5,
// 		Window: "1s",
// 	}
// 	req := httptest.NewRequest("GET", "http://nice.url?name=name&category=category", nil)
// 	req.RemoteAddr = "127.0.0.1"

// 	c, err := check.NewRateLimit(cfg)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}

// 	for i := 0; i < 5; i++ {
// 		_, err := c.Execute(context.Background(), req)
// 		if err != nil {
// 			t.Fatalf("expected no error, but got %v", err)
// 		}
// 	}

// 	time.Sleep(1100 * time.Millisecond)

// 	_, err = c.Execute(context.Background(), req)
// 	if err != nil {
// 		t.Fatalf("expected no error, but got %v", err)
// 	}
// }
