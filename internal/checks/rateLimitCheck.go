package checks

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const RateLimitCheckType = "rate_limit"

type RateLimitCheck struct {
	reqCount map[string]*ReqCount
	limit    int
	window   time.Duration
	mu       sync.Mutex
}

type ReqCount struct {
	count     int
	resetTime time.Time
}

func NewRateLimitCheck(cfg *config.RateLimitCheckConfig) (*RateLimitCheck, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	w, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, ErrInvalidWindow
	}
	if cfg.Limit <= 0 {
		return nil, ErrInvalidLimit
	}
	return &RateLimitCheck{
		reqCount: map[string]*ReqCount{},
		limit:    cfg.Limit,
		window:   w,
		mu:       sync.Mutex{},
	}, nil
}

func (c *RateLimitCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	count, ok := c.reqCount[r.RemoteAddr]
	if !ok {
		c.reqCount[r.RemoteAddr] = &ReqCount{
			count:     1,
			resetTime: now.Add(c.window),
		}
		return ctx, nil
	}

	if count.resetTime.Before(time.Now()) {
		count.count = 0
		count.resetTime = now.Add(c.window)
	}

	if count.count >= c.limit {
		return ctx, ErrTooManyRequests
	}

	count.count += 1
	return ctx, nil
}
