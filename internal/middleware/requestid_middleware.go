package middleware

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/ctxkeys"
	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var id string
		if v := r.Header.Get("X-Request-ID"); v == "" {
			id = uuid.New().String()
		} else {
			id = v
		}

		ctx := context.WithValue(r.Context(), ctxkeys.CtxRequestIDKey, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		w.Header().Set("X-Request-ID", id)

	})
}
