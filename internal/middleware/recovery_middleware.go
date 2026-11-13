package middleware

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/httperror"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				httperror.Write(w, http.StatusInternalServerError, httperror.CodeInternalError, "internal error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
