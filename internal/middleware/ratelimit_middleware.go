package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/ctxkeys"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

// TODO: По хорошему, надо другой способ идентифицировать клиента, например, тот же identity middleware или что-то лучше
func getClientID(r *http.Request) string {

	if v := r.Context().Value(ctxkeys.CtxUserIDKey); v != nil {
		return v.(string)
	}

	if v := r.Context().Value(ctxkeys.CtxRealIPKey); v != nil {
		return v.(string)
	}

	return r.RemoteAddr
}

type ReqCount struct {
	count     int
	resetTime time.Time
}

type RateLimit struct {
	mu     sync.Mutex
	rates  map[string]ReqCount
	limit  int
	window time.Duration
}

func NewRateLimit(limit int, window time.Duration) *RateLimit {
	return &RateLimit{
		mu:     sync.Mutex{},
		rates:  map[string]ReqCount{},
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimit) Middleware() pipeline.Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id := getClientID(r)

			rl.mu.Lock()

			now := time.Now()

			count, ok := rl.rates[id]

			if !ok || count.resetTime.Before(now) {
				count = ReqCount{
					count:     0,
					resetTime: now.Add(rl.window),
				}
			}

			if count.count >= rl.limit {
				rl.mu.Unlock()

				http.Error(
					w,
					"too many requests",
					http.StatusTooManyRequests,
				)
				return
			}

			count.count++

			rl.rates[id] = count

			rl.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
