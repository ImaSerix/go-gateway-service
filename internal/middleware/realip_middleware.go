package middleware

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/ctxkeys"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		if v := r.Header.Get("X-Forwarded-For"); v != "" {
			ip = v
		}
		if v := r.Header.Get("X-Real-IP"); v != "" {
			ip = v
		}

		ctx := context.WithValue(r.Context(), ctxkeys.CtxRealIPKey, ip)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
