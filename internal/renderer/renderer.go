package renderer

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Renderer interface {
	Render(s string, r *http.Request) (string, error)
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
