package check

import (
	"bytes"
	"context"
	"io"
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
	transforms     []pipeline.Transformer
	upstream       UpstreamClient
	store          Store
	expectedStatus int
}

func NewPolicyCheck(transforms []pipeline.Transformer, upstream UpstreamClient, store Store, expectedStatus int) *PolicyCheck {
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	return &PolicyCheck{
		transforms:     transforms,
		upstream:       upstream,
		store:          store,
		expectedStatus: expectedStatus,
	}
}

func (c *PolicyCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	newRequest, err := cloneRequestWithBody(r, ctx)
	if err != nil {
		return ctx, err
	}

	for _, transform := range c.transforms {
		err = transform.Transform(newRequest)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.upstream.Do(newRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.expectedStatus {
		return ctx, ErrUnauthorized
	}

	newCtx, err := c.store.Save(ctx, resp)
	if err != nil {
		return nil, err
	}

	return newCtx, nil
}

func cloneRequestWithBody(r *http.Request, ctx context.Context) (*http.Request, error) {
	clone := r.Clone(ctx)
	if r.Body == nil {
		return clone, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	clone.Body = io.NopCloser(bytes.NewBuffer(body))
	clone.ContentLength = int64(len(body))

	return clone, nil
}
