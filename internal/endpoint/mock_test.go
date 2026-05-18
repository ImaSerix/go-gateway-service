package endpoint_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type mockCheck struct {
	name   string
	err    error
	calls  *[]string
	setCtx bool
}

func (m *mockCheck) Execute(
	ctx context.Context,
	r *http.Request,
) (context.Context, error) {

	*m.calls = append(*m.calls, m.name)

	if m.err != nil {
		return ctx, m.err
	}

	if m.setCtx {
		ctx = context.WithValue(ctx, "user_id", 1001)
	}

	return ctx, nil
}

type mockTransformer struct {
	name      string
	err       error
	calls     *[]string
	expectCtx bool
}

func (m *mockTransformer) Transform(
	ctx context.Context,
	r *http.Request,
) error {

	*m.calls = append(*m.calls, m.name)

	if m.expectCtx {
		if v := ctx.Value("user_id"); v != 1001 {
			return errors.New("context value not propagated")
		}
	}

	return m.err
}

type mockProxy struct {
	calls *[]string
}

func (m *mockProxy) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	*m.calls = append(*m.calls, "proxy")
	w.WriteHeader(http.StatusOK)
}

type mockMiddleware struct {
	name  string
	calls *[]string
	pass  bool
}

func (m *mockMiddleware) Middleware() pipeline.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*m.calls = append(*m.calls, m.name)

			if !m.pass {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
