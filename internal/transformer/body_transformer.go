package transformer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ImaSerix/go-gateway-service/internal/resolver"
)

type BodyTransformer struct {
	template map[string]any
	resolver resolver.Resolver
}

func NewBodyTransformer(bodyBindings map[string]any, resolver resolver.Resolver) *BodyTransformer {
	return &BodyTransformer{
		template: bodyBindings,
		resolver: resolver,
	}
}

func mergeWithOverride(boundTemplate map[string]any, body any) any {

	m, ok := body.(map[string]any)
	if !ok {
		if boundTemplate == nil {
			return body
		}
		return boundTemplate
	}

	b := DeepCopy(m)

	for key, v := range boundTemplate {
		if _, ok := b[key]; !ok {
			b[key] = v
			continue
		}

		m, ok := v.(map[string]any)
		if mb, ok2 := b[key].(map[string]any); ok && ok2 {
			b[key] = mergeWithOverride(m, mb)
			continue
		}

		b[key] = v
	}

	return b
}

func isPlaceholder(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

// На данный момент работает только с map, в любом их виде, рекурсивно, не поддерживает списки
func (t *BodyTransformer) Bind(r *http.Request, template map[string]any) (map[string]any, error) {

	layer := DeepCopy(template)

	for key, v := range layer {

		if m, ok := v.(map[string]any); ok {
			bindedLayer, err := t.Bind(r, m)
			if err != nil {
				return nil, err
			}
			layer[key] = bindedLayer
			continue
		}

		if s, ok := v.(string); ok {
			if !isPlaceholder(s) {
				continue
			}
			v, ok := t.resolver.Resolve(r, s)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrInvalidKey, layer[key])
			}
			layer[key] = v
			continue
		}
	}

	return layer, nil
}

func (t *BodyTransformer) Transform(r *http.Request) error {

	if r == nil {
		return ErrNilRequest
	}

	if v := r.Header.Get("Content-Type"); !strings.HasPrefix(v, "application/json") {
		return ErrUnsupportedContentType
	}

	boundTemplate, err := t.Bind(r, t.template)
	if err != nil {
		return err
	}

	var data any
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("transform body: %w", err)
	}
	r.Body.Close()

	newBody := mergeWithOverride(boundTemplate, data)
	b, err := json.Marshal(newBody)
	if err != nil {
		return fmt.Errorf("transform body: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))

	return nil
}
