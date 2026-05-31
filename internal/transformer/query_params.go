package transformer

import (
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type QueryParams struct {
	queryBindings map[string]any
	render        renderer.Renderer
}

func NewQueryParams(qb map[string]any, render renderer.Renderer) *QueryParams {
	return &QueryParams{
		queryBindings: qb,
		render:        render,
	}
}

func (t *QueryParams) Transform(r *http.Request) error {

	if r == nil {
		return ErrNilRequest
	}

	q := r.URL.Query()

	for key, template := range t.queryBindings {
		if s, ok := template.(string); ok {
			v, err := t.render.Render(s, r)
			if err != nil {
				return fmt.Errorf("%w: %s", err, s)
			}

			q.Set(key, v)
			continue
		}
		q.Set(key, fmt.Sprint(template))
	}

	r.URL.RawQuery = q.Encode()

	return nil
}
