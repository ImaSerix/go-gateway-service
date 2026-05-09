package checks

import (
	"context"
	"net/http"
	"slices"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const IPWhiteListCheckType = "ip_whitelist"

type IPWhitelistCheck struct {
	ip []string
}

func NewIPWhiteListCheck(cfg *config.IPWhiteListCheckConfig) (*IPWhitelistCheck, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if len(cfg.IP) == 0 {
		return nil, ErrEmptyIP
	}

	return &IPWhitelistCheck{
		ip: cfg.IP,
	}, nil
}

func (c *IPWhitelistCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	if !slices.Contains(c.ip, r.RemoteAddr) {
		return ctx, ErrForbidden
	}

	return ctx, nil
}
