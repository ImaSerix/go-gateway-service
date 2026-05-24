package check

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type UpstreamClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type Store interface {
	Save(ctx context.Context, w *http.Response) (context.Context, error)
}

type PolicyCheck struct {
	transform      pipeline.Transformer
	upstream       UpstreamClient
	store          Store
	expectedStatus int
}

func NewPolicyCheck(transform pipeline.Transformer, upstream UpstreamClient, store Store) *PolicyCheck {
	return &PolicyCheck{
		transform: transform,
		upstream:  upstream,
		store:     store,
	}
}

func (c *PolicyCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {

	newRequest, err := http.NewRequest("", "", nil)
	if err != nil {
		return ctx, err
	}

	err = c.transform.Transform(newRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.upstream.Do(newRequest)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != c.expectedStatus {
		return ctx, ErrUnauthorized
	}

	newCtx, err := c.store.Save(ctx, resp)
	if err != nil {
		return nil, err
	}

	return newCtx, nil
}
