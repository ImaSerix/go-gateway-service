package renderer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Renderer interface {
	Render(s string, r *http.Request) (string, error)
	RenderURL(s string, r *http.Request) (string, error)
}

type ResponseRenderer interface {
	Render(s string, r *http.Response) (string, error)
}

var re = regexp.MustCompile(`\{([a-zA-Z0-9_]+):([a-zA-Z0-9_-]+)\}`)

type Resolver interface {
	Resolve(r *http.Request, source string, key string) (any, bool)
}

type ResponseResolver interface {
	Resolve(r *http.Response, source string, key string) (any, bool)
}

type Render struct {
	resolver Resolver
}

type ResponseRender struct {
	resolver ResponseResolver
}

func NewRender(resolver Resolver) *Render {
	return &Render{
		resolver: resolver,
	}
}

func NewResponseRender(resolver ResponseResolver) *ResponseRender {
	return &ResponseRender{
		resolver: resolver,
	}
}

func (rend *Render) Render(s string, r *http.Request) (string, error) {

	if strings.Count(s, "{") != strings.Count(s, "}") {
		return "", fmt.Errorf("%w: %s", ErrInvalidTemplate, "all brackets is not closed")
	}

	out, err := rend.render(s, r, nil)
	if err != nil {
		return "", err
	}

	return out, nil
}

// RenderURL renders a URL template so that a value stays a value: before "?" it is
// escaped as a path segment, after "?" as a query value, and cannot change the URL structure.
// A value that is "." or "..", or makes such a path segment, fails with ErrDotSegment;
// this error is returned even if another template of the URL cannot be resolved.
func (rend *Render) RenderURL(s string, r *http.Request) (string, error) {

	if strings.Count(s, "{") != strings.Count(s, "}") {
		return "", fmt.Errorf("%w: %s", ErrInvalidTemplate, "all brackets is not closed")
	}

	path, query, hasQuery := strings.Cut(s, "?")

	out, err := rend.render(path, r, func(v string) (string, error) {
		if v == "." || v == ".." {
			return "", ErrDotSegment
		}
		return url.PathEscape(v), nil
	})
	if errors.Is(err, ErrDotSegment) {
		return "", err
	}

	// An escaped value has no "/": the template and the result have the same segments.
	// Dot segments written in the template itself are left as they are.
	segments := strings.Split(path, "/")
	for i, segment := range strings.Split(out, "/") {
		if (segment == "." || segment == "..") && segment != segments[i] {
			return "", ErrDotSegment
		}
	}

	if err != nil {
		return "", err
	}

	if !hasQuery {
		return out, nil
	}

	q, err := rend.render(query, r, func(v string) (string, error) {
		return url.QueryEscape(v), nil
	})
	if err != nil {
		return "", err
	}

	return out + "?" + q, nil
}

// render replaces the templates in s; escape, if set, is applied to every value.
// On error it also returns s with the failed templates left as they are.
// The first error is returned, an error of escape goes before a failed resolve.
func (rend *Render) render(s string, r *http.Request, escape func(string) (string, error)) (string, error) {

	var renderErr error

	out := re.ReplaceAllStringFunc(s, func(s string) string {

		parts := re.FindStringSubmatch(s)

		source := parts[1]
		key := parts[2]

		v, ok := rend.resolver.Resolve(r, source, key)
		if !ok {
			if renderErr == nil {
				renderErr = fmt.Errorf("%w: %s", ErrFailedResolve, v)
			}
			return s
		}

		value := fmt.Sprint(v)
		if escape != nil {
			// chi routes by the escaped path when the request has one,
			// a route parameter is then escaped too: escape it only once.
			if source == "route" && r.URL.RawPath != "" {
				if unescaped, err := url.PathUnescape(value); err == nil {
					value = unescaped
				}
			}

			escaped, err := escape(value)
			if err != nil {
				if renderErr == nil || errors.Is(renderErr, ErrFailedResolve) {
					renderErr = err
				}
				return s
			}
			value = escaped
		}
		return value
	})

	return out, renderErr

}

func (rend *ResponseRender) Render(s string, r *http.Response) (string, error) {

	if strings.Count(s, "{") != strings.Count(s, "}") {
		return "", fmt.Errorf("%w: %s", ErrInvalidTemplate, "all brackets is not closed")
	}

	var renderErr error

	out := re.ReplaceAllStringFunc(s, func(s string) string {

		parts := re.FindStringSubmatch(s)

		source := parts[1]
		key := parts[2]

		v, ok := rend.resolver.Resolve(r, source, key)
		if !ok {
			renderErr = fmt.Errorf("%w: %s", ErrFailedResolve, v)
			return ""
		}
		return fmt.Sprint(v)
	})

	if renderErr != nil {
		return "", renderErr
	}

	return out, nil

}
