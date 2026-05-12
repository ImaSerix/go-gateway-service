package transformer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/transformer"
)

func TestBodyTransformer_Transform(t *testing.T) {
	tests := []struct {
		name      string
		template  map[string]any
		ctxValues map[string]any
		body      any
		expBody   any
		expErr    error
	}{
		{
			name: "success",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body: map[string]any{
				"user": map[string]any{
					"username": "usrname",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"username": "usrname",
					"id":       1001.0,
				},
			},
			expErr: nil,
		},
		{
			name: "nil body",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body: nil,
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001.0,
				},
			},
			expErr: nil,
		},
		{
			name: "override body value",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body: map[string]any{
				"user": map[string]any{
					"id": 1000,
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001.0,
				},
			},
			expErr: nil,
		},
		{
			name: "body primitive type",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body: "string",
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001.0,
				},
			},
			expErr: nil,
		},
		{
			name: "body list",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body: []any{"string", "string2"},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001.0,
				},
			},
			expErr: nil,
		},
		{
			name:     "bindings nil, body primitive",
			template: nil,
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body:    "string",
			expBody: "string",
			expErr:  nil,
		},
		{
			name:     "bindings nil, body list",
			template: nil,
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			body:    []any{"string", "string2"},
			expBody: []any{"string", "string2"},
			expErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			templateCopy := transformer.DeepCopy(test.template)

			tr := transformer.NewBodyTransformer(test.template)

			b, err := json.Marshal(test.body)
			if err != nil {
				t.Fatalf("expected no error on marshal, but got %v", err)
			}
			r := httptest.NewRequest("GET", "http://nice.url", nil)
			r.Body = io.NopCloser(bytes.NewBuffer(b))
			r.ContentLength = int64(len(b))
			r.Header.Set("Content-Type", "application/json")

			ctx := context.Background()
			for ctxKey, v := range test.ctxValues {
				ctx = context.WithValue(ctx, ctxKey, v)
			}

			err = tr.Transform(ctx, r)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, but got %v", test.expErr, err)
			}

			var body any
			json.NewDecoder(r.Body).Decode(&body)
			if !reflect.DeepEqual(body, test.expBody) {
				t.Fatalf("expected body %v, but got %v", test.expBody, body)
			}

			if !reflect.DeepEqual(test.template, templateCopy) {
				t.Fatal("expected template inmutable, but it has changed")
			}
		})
	}
}

func TestBodyTransformer_Transform_NilRequest(t *testing.T) {
	tr := transformer.NewBodyTransformer(nil)

	err := tr.Transform(context.Background(), nil)
	if !errors.Is(err, transformer.ErrNilRequest) {
		t.Fatalf("expected erorr %v, but got %v", transformer.ErrNilRequest, err)
	}
}

func TestBodyTransformer_Transform_UnsupportedContentType(t *testing.T) {
	tr := transformer.NewBodyTransformer(nil)

	r := httptest.NewRequest("GET", "http://nice.url", nil)

	err := tr.Transform(context.Background(), r)
	if !errors.Is(err, transformer.ErrUnsupportedContentType) {
		t.Fatalf("expected erorr %v, but got %v", transformer.ErrUnsupportedContentType, err)
	}
}

func TestBodyTransformer_Transform_NoKeyInContext(t *testing.T) {
	tr := transformer.NewBodyTransformer(map[string]any{
		"user": map[string]any{
			"id": "{user_id}",
		},
	})

	r := httptest.NewRequest("GET", "http://nice.url", nil)
	r.Header.Set("Content-Type", "application/json")

	err := tr.Transform(context.Background(), r)
	if !errors.Is(err, transformer.ErrNoKeyInContext) {
		t.Fatalf("expected erorr %v, but got %v", transformer.ErrNoKeyInContext, err)
	}
}

func TestBodyTransformer_Bind(t *testing.T) {

	tests := []struct {
		name      string
		template  map[string]any
		ctxValues map[string]any
		expBody   map[string]any
		expErr    error
	}{
		{
			name: "success",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001,
				},
			},
			ctxValues: map[string]any{
				"user_id": 1001,
			},
			expErr: nil,
		},
		{
			name: "more than one",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
				"chat": map[string]any{
					"id":      "{chat_id}",
					"content": "{chat_content}",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001,
				},
				"chat": map[string]any{
					"id":      1002,
					"content": "nice content",
				},
			},
			ctxValues: map[string]any{
				"user_id":      1001,
				"chat_id":      1002,
				"chat_content": "nice content",
			},
			expErr: nil,
		},
		{
			name: "no key in context",
			template: map[string]any{
				"user": map[string]any{
					"id": "{user_id}",
				},
				"chat": map[string]any{
					"id":      "{chat_id}",
					"content": "{chat_content}",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001,
				},
				"chat": map[string]any{
					"id":      1002,
					"content": "nice content",
				},
			},
			ctxValues: map[string]any{
				"chat_id":      1002,
				"chat_content": "nice content",
			},
			expErr: transformer.ErrNoKeyInContext,
		},
		{
			name: "placeholder not a string",
			template: map[string]any{
				"user": map[string]any{
					"id": 1001,
				},
				"chat": map[string]any{
					"id":      "{chat_id}",
					"content": "{chat_content}",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1001,
				},
				"chat": map[string]any{
					"id":      1002,
					"content": "nice content",
				},
			},
			ctxValues: map[string]any{
				"chat_id":      1002,
				"chat_content": "nice content",
			},
			expErr: nil,
		},
		{
			name: "invalid placeholder",
			template: map[string]any{
				"user": map[string]any{
					"id": "{id",
				},
				"chat": map[string]any{
					"id":      "{chat_id}",
					"content": "{chat_content}",
				},
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": "{id",
				},
				"chat": map[string]any{
					"id":      1002,
					"content": "nice content",
				},
			},
			ctxValues: map[string]any{
				"chat_id":      1002,
				"chat_content": "nice content",
			},
			expErr: nil,
		},
		{
			name: "template without placeholders",
			template: map[string]any{
				"user": map[string]any{
					"id": 1000,
				},
				"chat": map[string]any{
					"id":      2002,
					"content": "bad content",
				},
				"ok": true,
			},
			expBody: map[string]any{
				"user": map[string]any{
					"id": 1000,
				},
				"chat": map[string]any{
					"id":      2002,
					"content": "bad content",
				},
				"ok": true,
			},
			ctxValues: map[string]any{
				"chat_id":      1002,
				"chat_content": "nice content",
			},
			expErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctx := context.Background()

			for key, v := range test.ctxValues {
				ctx = context.WithValue(ctx, key, v)
			}

			bt := transformer.NewBodyTransformer(test.template)

			binded, err := bt.Bind(ctx, test.template)
			if !errors.Is(err, test.expErr) {
				t.Fatalf("expected error %v, but got %v", test.expErr, err)
			}

			if err == nil && !reflect.DeepEqual(binded, test.expBody) {
				t.Fatalf("expected body %v, but got %v", test.expBody, test.template)
			}
		})
	}

}
