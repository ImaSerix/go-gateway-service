package middleware

import (
	"net/http"
	"strings"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

func CORS(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
) pipeline.Middleware {

	allowed := make(map[string]struct{}, len(allowedOrigins))
	allowAll := false

	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAll = true
			break
		}

		allowed[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			corsAllowed := false
			if allowAll {
				corsAllowed = true

				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if _, ok := allowed[origin]; ok {

				corsAllowed = true

				w.Header().Set(
					"Access-Control-Allow-Origin",
					origin,
				)

				w.Header().Set("Vary", "Origin")
			}

			if corsAllowed {
				if len(allowedMethods) != 0 {
					w.Header().Set(
						"Access-Control-Allow-Methods",
						strings.Join(allowedMethods, ", "),
					)
				}

				if len(allowedHeaders) != 0 {
					w.Header().Set(
						"Access-Control-Allow-Headers",
						strings.Join(allowedHeaders, ", "),
					)
				}
			}

			if r.Method == http.MethodOptions && corsAllowed {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
