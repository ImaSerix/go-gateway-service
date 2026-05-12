package transformer

import (
	"context"
	"fmt"
	"net/http"
)

type HeaderTransformer struct {
	headerBindings map[string]string
}

func NewHeaderTransformer(headerBindings map[string]string) *HeaderTransformer {
	return &HeaderTransformer{
		headerBindings: headerBindings,
	}
}

func (t *HeaderTransformer) Transform(ctx context.Context, r *http.Request) error {

	if r == nil {
		return ErrNilRequest
	}

	for header, ctxKey := range t.headerBindings {
		v := ctx.Value(ctxKey)
		if v == nil {
			return fmt.Errorf("%w: %s", ErrNoKeyInContext, ctxKey)
		}

		r.Header.Set(header, fmt.Sprint(v))
	}

	return nil
}
