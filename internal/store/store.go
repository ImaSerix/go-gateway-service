package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

type Store struct {
	template map[string]string
	render   renderer.Render
}

func NewStore(template map[string]string, render renderer.Render) *Store {
	return &Store{
		template: template,
		render:   render,
	}
}

func (s *Store) Save(ctx context.Context, req *http.Request) (context.Context, error) {

	if req == nil {
		return ctx, ErrNilRequest
	}

	for k, t := range s.template {

		v, err := s.render.Render(t, req)
		if err != nil {
			return ctx, fmt.Errorf("store save: %w", err)
		}

		ctx = context.WithValue(ctx, k, v)
	}

	return ctx, nil
}
