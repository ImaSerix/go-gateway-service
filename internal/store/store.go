package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Store struct {
	template map[string]string
	render   renderer.ResponseRenderer
}

func NewStore(template map[string]string, render renderer.ResponseRenderer) *Store {
	return &Store{
		template: template,
		render:   render,
	}
}

func (s *Store) Save(ctx context.Context, res *http.Response) (context.Context, error) {

	if res == nil {
		return ctx, ErrNilResponse
	}

	for k, t := range s.template {

		v, err := s.render.Render(t, res)
		if err != nil {
			return ctx, fmt.Errorf("store save: %w", err)
		}

		ctx = context.WithValue(ctx, k, v)
	}

	return ctx, nil
}
