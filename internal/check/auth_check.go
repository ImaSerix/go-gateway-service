package check

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/ImaSerix/go-gateway-service/internal/config"
// )

// const AuthCheckType = "auth"

// type Auth struct {
// 	url            string
// 	method         Method
// 	forwardHeaders map[string]string
// 	storeBody      map[string]string
// 	storeHeaders   map[string]string
// 	expectedStatus int
// 	client         *http.Client
// }

// func NewAuth(cfg config.AuthCheck, client *http.Client) (*Auth, error) {

// 	if cfg.URL == "" {
// 		return nil, ErrEmptyURL
// 	}

// 	url, err := url.Parse(cfg.URL)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, cfg.URL)
// 	}
// 	if url.Host == "" {
// 		return nil, ErrEmptyHost
// 	}

// 	m, err := NewMethod(cfg.Method)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if client == nil {
// 		return nil, ErrNilClient
// 	}

// 	expectedStatus := cfg.ExpectedStatus
// 	if expectedStatus == 0 {
// 		expectedStatus = http.StatusOK
// 	}

// 	if expectedStatus < 100 || expectedStatus > 599 {
// 		return nil, ErrInvalidExpectedStatus
// 	}

// 	return &Auth{
// 		url:            cfg.URL,
// 		method:         m,
// 		forwardHeaders: cfg.ForwardHeaders,
// 		storeBody:      cfg.Store.Body,
// 		storeHeaders:   cfg.Store.Headers,
// 		expectedStatus: expectedStatus,
// 		client:         client,
// 	}, nil
// }

// func (c *Auth) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
// 	if r == nil {
// 		return ctx, ErrNilRequest
// 	}

// 	req, err := http.NewRequestWithContext(ctx, string(c.method), c.url, nil)
// 	if err != nil {
// 		return ctx, fmt.Errorf("auth check: create request: %w", err)
// 	}

// 	for newHeader, oldHeader := range c.forwardHeaders {
// 		if val := r.Header.Get(oldHeader); val != "" {
// 			req.Header.Set(newHeader, val)
// 		}
// 	}

// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return ctx, fmt.Errorf("auth check: request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != c.expectedStatus {
// 		return ctx, ErrUnauthorized
// 	}

// 	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
// 		var body any
// 		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil && err != io.EOF {
// 			return ctx, fmt.Errorf("auth check: decode body: %w", err)
// 		}

// 		if m, ok := body.(map[string]any); ok {
// 			for ctxKey, bodyKey := range c.storeBody {
// 				if val, ok := m[bodyKey]; ok {
// 					ctx = context.WithValue(ctx, ctxKey, val)
// 				}
// 			}
// 		}
// 	}

// 	for ctxKey, headerKey := range c.storeHeaders {
// 		if val := resp.Header.Get(headerKey); val != "" {
// 			ctx = context.WithValue(ctx, ctxKey, val)
// 		}
// 	}

// 	return ctx, nil
// }
