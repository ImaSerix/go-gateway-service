package checks

import (
	"context"
	"net/http"
	"slices"
	"strings"

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

	split := strings.Split(r.RemoteAddr, ":")

	if !slices.Contains(c.ip, split[0]) {
		return ctx, ErrForbidden
	}

	return ctx, nil
}
