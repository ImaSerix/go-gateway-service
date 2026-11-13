package middleware

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

// StripHeaders removes the listed headers from the incoming request,
// so a value sent by the client never reaches checks, transforms or upstreams.
func StripHeaders(headers []string) pipeline.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			for _, h := range headers {
				r.Header.Del(h)
			}

			next.ServeHTTP(w, r)
		})
	}
}
