package transformer_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/transformer"
)

func TestQueryParams_Transform_Success(t *testing.T) {

	hb := map[string]any{
		"country": "{ctx.country}",
		"age":     "{ctx.age}",
	}

	render := &renderMock{
		values: map[string]string{
			"{ctx.country}": "Latvia",
			"{ctx.age}":     "21",
		},
	}

	tf := transformer.NewQueryParams(hb, render)

	r := httptest.NewRequest("GET", "http://nice.url", nil)

	err := tf.Transform(r)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if v := r.URL.Query().Get("country"); v != "Latvia" {
		t.Fatalf("expected get param country 'Latvia', but got %s", v)
	}

	if v := r.URL.Query().Get("age"); v != "21" {
		t.Fatalf("expected get param age '21', but got %s", v)
	}

}

func TestQueryParams_Transform_NilRequest(t *testing.T) {

	hb := map[string]any{
		"country": "{ctx.country}",
		"age":     "{ctx.age}",
	}

	ctx := context.WithValue(t.Context(), "user_username", "nice username")
	ctx = context.WithValue(ctx, "user_id", 1001)

	tf := transformer.NewQueryParams(hb, nil)

	err := tf.Transform(nil)
	if !errors.Is(err, transformer.ErrNilRequest) {
		t.Fatalf("expected wrapped error %v, but got %v", transformer.ErrNilRequest, err)
	}

}

func TestQueryParams_Transform_KeyNotInContext(t *testing.T) {

	hb := map[string]any{
		"country": "{ctx.country}",
		"age":     "{ctx.age}",
	}

	badError := errors.New("bad error")

	render := &renderMock{
		values: map[string]string{},
		err:    badError,
	}

	tf := transformer.NewQueryParams(hb, render)

	r := httptest.NewRequest("GET", "http://nice.url", nil)

	err := tf.Transform(r)
	if !errors.Is(err, badError) {
		t.Fatalf("expected wrapped error %v, but got %v", badError, err)
	}

}
