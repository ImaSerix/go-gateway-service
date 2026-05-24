package transformer

import (
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Headers struct {
	headerBindings map[string]string
	render         renderer.Renderer
}

func NewHeaderTransformer(headerBindings map[string]string, render renderer.Renderer) *Headers {
	return &Headers{
		headerBindings: headerBindings,
		render:         render,
	}
}

func (t *Headers) Transform(r *http.Request) error {

	if r == nil {
		return ErrNilRequest
	}

	for header, key := range t.headerBindings {
		v, err := t.render.Render(key, r)
		if err != nil {
			return fmt.Errorf("%w: %s", err, key)
		}

		r.Header.Set(header, fmt.Sprint(v))
	}

	return nil
}
