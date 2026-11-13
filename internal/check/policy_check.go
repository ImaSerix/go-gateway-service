package check

import (
	"bytes"
	"context"
	"fmt"
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
	forwardBody    bool
}

func NewPolicyCheck(transforms []pipeline.Transformer, upstream UpstreamClient, store Store, expectedStatus int, forwardBody bool) *PolicyCheck {
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	return &PolicyCheck{
		transforms:     transforms,
		upstream:       upstream,
		store:          store,
		expectedStatus: expectedStatus,
		forwardBody:    forwardBody,
	}
}

func (c *PolicyCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return ctx, ErrNilRequest
	}

	var newRequest *http.Request
	var err error

	if c.forwardBody {
		newRequest, err = cloneRequestWithBody(r, ctx)
		if err != nil {
			return ctx, err
		}
	} else {
		// The client body is neither read nor sent to the policy service.
		newRequest = r.Clone(ctx)
		newRequest.Body = http.NoBody
		newRequest.ContentLength = 0
		newRequest.TransferEncoding = nil
		newRequest.Header.Del("Content-Length")
		newRequest.Header.Del("Content-Type")
	}

	for _, transform := range c.transforms {
		err = transform.Transform(newRequest)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.upstream.Do(newRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPolicyUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.expectedStatus {
		if resp.StatusCode == http.StatusUnauthorized {
			return ctx, ErrUnauthorized
		}

		return ctx, fmt.Errorf("%w: status %d", ErrPolicyUnavailable, resp.StatusCode)
	}

	newCtx, err := c.store.Save(ctx, resp)
	if err != nil {
		// The policy service answered with the expected status, but its response is unusable.
		return nil, fmt.Errorf("%w: %v", ErrPolicyUnavailable, err)
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
