package checks_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"gopkg.in/yaml.v3"
)

func TestCheck_AuthCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.AuthCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.AuthCheckConfig{
		URL: "http://nice.url",
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.AuthCheck); !ok {
		t.Fatal("expected check with type 'checks.AuthCheck', but got something other")
	}
}

func TestCheck_HeaderRequiredCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.HeaderRequiredCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.HeaderRequiredCheckConfig{
		Headers: []string{
			"X-Username",
			"X-Password",
		},
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.HeaderRequiredCheck); !ok {
		t.Fatal("expected check with type 'checks.HeaderRequiredCheck', but got something other")
	}
}

func TestCheck_InjectCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.InjectCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.InjectCheckConfig{
		Ctx: map[string]any{
			"key":  "value",
			"key2": 21,
		},
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.InjectCheck); !ok {
		t.Fatal("expected check with type 'checks.InjectCheck', but got something other")
	}
}

func TestCheck_IPWhitelistCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.IPWhiteListCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.IPWhiteListCheckConfig{
		IP: []string{
			"127.0.0.1",
			"192.168.0.1",
		},
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.IPWhitelistCheck); !ok {
		t.Fatal("expected check with type 'checks.IPWhitelistCheck', but got something other")
	}
}

func TestCheck_TimeoutCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.TimeoutCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.TimeoutCheckConfig{
		Duration: "2s",
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.TimeoutCheck); !ok {
		t.Fatal("expected check with type 'checks.TimeoutCheck', but got something other")
	}
}

func TestCheck_QueryRequiredCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.QueryRequiredCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.QueryRequiredCheckConfig{
		Queries: []string{
			"name",
			"category",
		},
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.QueryRequiredCheck); !ok {
		t.Fatal("expected check with type 'checks.QueryRequiredCheck', but got something other")
	}
}

func TestCheck_RateLimitCheck(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   checks.RateLimitCheckType,
		Config: yaml.Node{},
	}

	checkCfg := &config.RateLimitCheckConfig{
		Limit:  50,
		Window: "1m",
	}

	cfg.Config.Encode(checkCfg)

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil check, but got nil")
	}
	if _, ok := c.(*checks.RateLimitCheck); !ok {
		t.Fatal("expected check with type 'checks.RateLimitCheck', but got something other")
	}
}

func TestCheck_UnsupportedType(t *testing.T) {
	cfg := config.CheckConfig{
		Type:   "unsupported",
		Config: yaml.Node{},
	}

	c, err := checks.CheckFactory(cfg, http.DefaultClient)
	if !errors.Is(err, checks.ErrUnsupportedType) {
		t.Fatalf("expected wrapped error %v, but got %v", checks.ErrUnsupportedType, err)
	}
	if c != nil {
		t.Fatalf("expected nil check, but got non-nil")
	}
}
