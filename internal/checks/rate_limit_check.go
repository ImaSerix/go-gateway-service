package checks

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const RateLimitCheckType = "rate_limit"

type RateLimitCheck struct {
	reqCount map[string]*ReqCount
	limit    int
	window   time.Duration
	mu       *sync.Mutex
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

	if w <= 0 {
		return nil, ErrInvalidWindow
	}

	if cfg.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	reqCount := map[string]*ReqCount{}
	mu := sync.Mutex{}

	//TODO: Надо делать какую-то clean-up функцию, чтоб это останавливали
	go func() {
		t := time.NewTicker(w)
		for {
			<-t.C
			mu.Lock()
			for k, v := range reqCount {
				if time.Now().Sub(v.resetTime) > w {
					delete(reqCount, k)
				}
			}
			mu.Unlock()
		}
	}()

	return &RateLimitCheck{
		reqCount: reqCount,
		limit:    cfg.Limit,
		window:   w,
		mu:       &mu,
	}, nil
}

// TODO: имеет смысл поменять способ проверки времени, возможно как-то инжектить time.Now() функцию, для более удобного теста
func (c *RateLimitCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	split := strings.Split(r.RemoteAddr, ":")

	count, ok := c.reqCount[split[0]]
	if !ok {
		c.reqCount[split[0]] = &ReqCount{
			count:     1,
			resetTime: now.Add(c.window),
		}
		return ctx, nil
	}

	if count.resetTime.Before(now) {
		count.count = 0
		count.resetTime = now.Add(c.window)

	}

	if count.count >= c.limit {
		return ctx, ErrTooManyRequests
	}

	count.count += 1
	return ctx, nil
}
