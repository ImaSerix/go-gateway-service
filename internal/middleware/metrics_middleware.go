package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

// TODO: В будущем поменять на что-то более внятное для метрик, а именно куда-то отсылать
type Metric struct {
	mu         sync.Mutex
	request    map[string]int
	latency    map[string]time.Duration
	statusCode map[string]int
}

func NewMetric() *Metric {
	return &Metric{
		mu:         sync.Mutex{},
		request:    map[string]int{},
		latency:    map[string]time.Duration{},
		statusCode: map[string]int{},
	}
}

func (m *Metric) Middleware() pipeline.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			sw := &statusWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(sw, r)

			duration := time.Since(start)
			key := r.Method + " " + r.URL.Path

			m.mu.Lock()
			m.request[key]++
			m.latency[key] += duration
			m.statusCode[key+" "+strconv.Itoa(sw.statusCode)]++
			m.mu.Unlock()

		})
	}
}
