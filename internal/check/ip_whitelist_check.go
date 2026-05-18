package check

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type IPWhiteList struct {
	ip []string
}

func NewIPWhiteList(cfg config.IPWhiteListCheck) (*IPWhiteList, error) {

	if len(cfg.IP) == 0 {
		return nil, ErrEmptyIP
	}

	return &IPWhiteList{
		ip: cfg.IP,
	}, nil
}

func (c *IPWhiteList) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	split := strings.Split(r.RemoteAddr, ":")

	if !slices.Contains(c.ip, split[0]) {
		return ctx, ErrForbidden
	}

	return ctx, nil
}
