package transformer_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/transformer"
)

func TestHeaderTransformer_Transform_Success(t *testing.T) {

	hb := map[string]string{
		"X-Username": "{ctx.user_username}",
		"X-User-ID":  "{ctx.user_id}",
	}

	render := &renderMock{
		values: map[string]string{
			"{ctx.user_username}": "nice username",
			"{ctx.user_id}":       "1001",
		},
	}

	tf := transformer.NewHeaderTransformer(hb, render)

	r := httptest.NewRequest("GET", "http://nice.url", nil)

	err := tf.Transform(r)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if r.Header.Get("X-Username") != "nice username" {
		t.Fatal("expected header 'X-Username' be 'nice username', but got something other")
	}
	if r.Header.Get("X-User-ID") != "1001" {
		t.Fatal("expected header 'X-User-ID' be '1001', but got something other")
	}

}

func TestHeaderTransformer_Transform_NilRequest(t *testing.T) {

	hb := map[string]string{
		"X-Username": "user_username",
		"X-User-ID":  "user_id",
	}

	ctx := context.WithValue(t.Context(), "user_username", "nice username")
	ctx = context.WithValue(ctx, "user_id", 1001)

	tf := transformer.NewHeaderTransformer(hb, nil)

	err := tf.Transform(nil)
	if !errors.Is(err, transformer.ErrNilRequest) {
		t.Fatalf("expected wrapped error %v, but got %v", transformer.ErrNilRequest, err)
	}

}

func TestHeaderTransformer_Transform_KeyNotInContext(t *testing.T) {

	hb := map[string]string{
		"X-Username": "{ctx.user_username}",
		"X-User-ID":  "{ctx.user_id}",
	}

	badError := errors.New("bad error")

	render := &renderMock{
		values: map[string]string{
			"{ctx.user_id}": "1001",
		},
		err: badError,
	}

	tf := transformer.NewHeaderTransformer(hb, render)

	r := httptest.NewRequest("GET", "http://nice.url", nil)

	err := tf.Transform(r)
	if !errors.Is(err, badError) {
		t.Fatalf("expected wrapped error %v, but got %v", badError, err)
	}

}
