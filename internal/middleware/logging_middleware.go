package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		sw := &statusWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(sw, r)

		//TODO: Надо будет поменять это на какой-т кастомный логер, или прикритить интерфейс в целом
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "status code", sw.statusCode, "duration", time.Since(start))
	})
}
