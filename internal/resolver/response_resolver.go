package resolver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ResponseMultiResolver struct {
	resolvers map[string]ResponseResolver
}

type ResponseResolver interface {
	Resolve(r *http.Response, key string) (any, bool)
}

func NewResponseMultiResolver() *ResponseMultiResolver {
	r := &ResponseMultiResolver{
		resolvers: map[string]ResponseResolver{},
	}

	r.Register("header", NewResponseHeaderResolver())
	r.Register("body", NewResponseBodyResolver())

	return r
}

func (r *ResponseMultiResolver) Register(source string, resolver ResponseResolver) {
	r.resolvers[source] = resolver
}

func (r *ResponseMultiResolver) Resolve(res *http.Response, source string, key string) (any, bool) {
	resolver, ok := r.resolvers[source]
	if !ok {
		return nil, false
	}

	return resolver.Resolve(res, key)
}

type ResponseHeaderResolver struct{}

func NewResponseHeaderResolver() *ResponseHeaderResolver {
	return &ResponseHeaderResolver{}
}

func (r *ResponseHeaderResolver) Resolve(res *http.Response, key string) (any, bool) {
	if res == nil {
		return nil, false
	}

	v := res.Header.Get(key)
	if v == "" {
		return nil, false
	}

	return v, true
}

type ResponseBodyResolver struct{}

func NewResponseBodyResolver() *ResponseBodyResolver {
	return &ResponseBodyResolver{}
}

func (r *ResponseBodyResolver) Resolve(res *http.Response, key string) (any, bool) {
	if res == nil || res.Body == nil {
		return nil, false
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, false
	}
	res.Body.Close()
	res.Body = io.NopCloser(bytes.NewBuffer(body))

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, false
	}

	v, ok := data[key]
	return v, ok
}
