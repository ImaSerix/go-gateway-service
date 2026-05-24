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

var re = regexp.MustCompile(`\{([a-zA-Z0-9_]+):([a-zA-Z0-9_]+)\}`)

type Resolver interface {
	Resolve(r *http.Request, source string, key string) (any, bool)
}

type Render struct {
	resolver Resolver
}

func NewRender(resolver Resolver) *Render {
	return &Render{
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
