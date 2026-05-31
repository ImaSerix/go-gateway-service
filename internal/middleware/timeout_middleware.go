package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

func Timeout(timeout time.Duration) pipeline.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

		})
	}
}
