package checks_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/checks"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestAuthCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.AuthCheckConfig
		client *http.Client
		expErr error
	}{
		{
			name: "success",
			cfg: &config.AuthCheckConfig{
				URL: "http://nice.url",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method: "POST",
			},
			client: http.DefaultClient,
			expErr: nil,
		},
		{
			name: "invalid expected status",
			cfg: &config.AuthCheckConfig{
				URL: "http://nice.url",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method:         "POST",
				ExpectedStatus: 99,
			},
			client: http.DefaultClient,
			expErr: checks.ErrInvalidExpectedStatus,
		},
		{
			name:   "nil cfg",
			cfg:    nil,
			expErr: checks.ErrNilConfig,
		},
		{
			name: "invalid url",
			cfg: &config.AuthCheckConfig{
				URL: "://bad.url",
				ForwardHeaders: map[string]string{
					"X-Username":  "X-Username",
					"tokenHeader": "X-Token",
					"X-Password":  "X-Password",
				},
				Method: "POST",
			},
			client: http.DefaultClient,
			expErr: checks.ErrInvalidURL,
		},
		{
			name: "empty url",
			cfg: &config.AuthCheckConfig{
				URL: "",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method: "POST",
			},
			client: http.DefaultClient,
			expErr: checks.ErrEmptyURL,
		},
		{
			name: "empty host",
			cfg: &config.AuthCheckConfig{
				URL: "http://",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method: "POST",
			},
			client: http.DefaultClient,
			expErr: checks.ErrEmptyHost,
		},
		{
			name: "invalid method",
			cfg: &config.AuthCheckConfig{
				URL: "http://nice.url",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method: "INVALID",
			},
			client: http.DefaultClient,
			expErr: checks.ErrInvalidMethod,
		},
		{
			name: "nil client",
			cfg: &config.AuthCheckConfig{
				URL: "http://nice.url",
				ForwardHeaders: map[string]string{
					"X-Username": "X-Username",
					"X-Password": "X-Password",
				},
				Method: "GET",
			},
			client: nil,
			expErr: checks.ErrNilClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := checks.NewAuthCheck(test.cfg, test.client)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, but got %v", test.expErr, err)
			}
			if err == nil && c == nil {
				t.Fatal("got nil check, but no error")
			}
			if err != nil && c != nil {
				t.Fatal("got error, but check is not nil")
			}
		})
	}
}

func SetupAuthCheck_Execute(token string, client *http.Client, cfg *config.AuthCheckConfig) (*checks.AuthCheck, *httptest.Server, *http.Request, context.Context, *string, *string) {

	var gotUsername string
	var gotPassword string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		gotUsername = r.Header.Get("X-Username")
		gotPassword = r.Header.Get("X-Password")

		w.Header().Set("X-Token", token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	}))
	cfg.URL = srv.URL

	req := httptest.NewRequest("POST", "http://nice.target", nil)
	ctx := context.Background()

	req.Header.Set("X-Username", "nice username")
	req.Header.Set("X-Password", "nice password")

	c, _ := checks.NewAuthCheck(cfg, client)

	return c, srv, req, ctx, &gotUsername, &gotPassword
}
func TestAuthCheck_Execute_Success(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 200,
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	c, srv, req, ctx, gotUsername, gotPassword := SetupAuthCheck_Execute(token, http.DefaultClient, cfg)
	defer srv.Close()

	newCtx, err := c.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}
	if v := newCtx.Value("token"); v != token {
		t.Fatalf("expected in context token %s, but got %s", token, v)
	}
	if v := newCtx.Value("tokenHeader"); v != token {
		t.Fatalf("expected in context tokenHeader %s, but got %s", token, v)
	}
	if *gotUsername != "nice username" {
		t.Fatalf("expected username 'nice username', but got %s", *gotUsername)
	}
	if *gotPassword != "nice password" {
		t.Fatalf("expected password 'nice password', but got %s", *gotPassword)
	}
}

func TestAuthCheck_Execute_NilRequest(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 200,
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	c, srv, _, ctx, _, _ := SetupAuthCheck_Execute(token, nil, cfg)
	defer srv.Close()

	newCtx, err := c.Execute(ctx, nil)
	if !errors.Is(err, checks.ErrNilRequest) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
	if v := newCtx.Value("token"); v != nil {
		t.Fatalf("expected in context empty token, but got %s", v)
	}
	if v := newCtx.Value("tokenHeader"); v != nil {
		t.Fatalf("expected in context empty tokenHeader, but got %s", v)
	}
}

func TestAuthCheck_Execute_NilClient(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 200,
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	c, srv, req, ctx, _, _ := SetupAuthCheck_Execute(token, http.DefaultClient, cfg)
	srv.Close()

	newCtx, err := c.Execute(ctx, req)
	if !strings.Contains(err.Error(), "auth check: request failed:") {
		t.Fatalf("expected error contains 'auth check: request failed', but got error without %v", err)
	}
	if v := newCtx.Value("token"); v != nil {
		t.Fatalf("expected in context empty token, but got %s", v)
	}
	if v := newCtx.Value("tokenHeader"); v != nil {
		t.Fatalf("expected in context empty tokenHeader, but got %s", v)
	}
}

func TestAuthCheck_Execute_Unauthorized(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 300,
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	c, srv, req, ctx, _, _ := SetupAuthCheck_Execute(token, http.DefaultClient, cfg)
	defer srv.Close()

	newCtx, err := c.Execute(ctx, req)
	if !errors.Is(err, checks.ErrUnauthorized) {
		t.Fatalf("expected error %v, but got %v", checks.ErrNilRequest, err)
	}
	if v := newCtx.Value("token"); v != nil {
		t.Fatalf("expected in context empty token, but got %s", v)
	}
	if v := newCtx.Value("tokenHeader"); v != nil {
		t.Fatalf("expected in context empty tokenHeader, but got %s", v)
	}
}

func TestAuthCheck_Execute_Success_EmptyBody(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 200,
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	_, srv, req, ctx, _, _ := SetupAuthCheck_Execute(token, http.DefaultClient, cfg)
	srv.Close()

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	cfg.URL = srv.URL
	c, _ := checks.NewAuthCheck(cfg, http.DefaultClient)

	_, err := c.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}
}

func TestAuthCheck_Execute_Success_BodyPrimitive(t *testing.T) {
	cfg := &config.AuthCheckConfig{
		ForwardHeaders: map[string]string{
			"X-Username": "X-Username",
			"X-Password": "X-Password",
		},
		Method: "POST",
		Store: config.Store{
			Body: map[string]string{
				"token": "token",
			},
			Headers: map[string]string{
				"tokenHeader": "X-Token",
			},
		},
		ExpectedStatus: 200,
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	_, srv, req, ctx, _, _ := SetupAuthCheck_Execute(token, http.DefaultClient, cfg)
	srv.Close()

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("priimitive"))
	}))

	cfg.URL = srv.URL
	c, _ := checks.NewAuthCheck(cfg, http.DefaultClient)

	_, err := c.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}
}
