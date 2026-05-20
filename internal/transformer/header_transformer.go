package transformer

import (
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

type HeaderTransformer struct {
	headerBindings map[string]string
	resolver       resolver.Resolver
}

func NewHeaderTransformer(headerBindings map[string]string, resolver resolver.Resolver) *HeaderTransformer {
	return &HeaderTransformer{
		headerBindings: headerBindings,
		resolver:       resolver,
	}
}

func (t *HeaderTransformer) Transform(r *http.Request) error {

	if r == nil {
		return ErrNilRequest
	}

	for header, key := range t.headerBindings {
		v, ok := t.resolver.Resolve(r, key)
		if !ok {
			return fmt.Errorf("%w: %s", ErrInvalidKey, key)
		}

		r.Header.Set(header, fmt.Sprint(v))
	}

	return nil
}
