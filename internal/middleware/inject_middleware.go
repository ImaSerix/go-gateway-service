package middleware

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

func Inject(injectValues map[string]any) pipeline.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			for k, v := range injectValues {
				ctx = context.WithValue(ctx, k, v)
			}

			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
