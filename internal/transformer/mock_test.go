package transformer_test

import (
	"net/http"
)

type renderMock struct {
	values        map[string]string
	err           error
	forHeaderTest bool
}

func (rm *renderMock) Render(s string, r *http.Request) (string, error) {

	v, ok := rm.values[s]
	if !ok {
		return s, rm.err
	}
	return v, rm.err
}

func (rm *renderMock) RenderURL(s string, r *http.Request) (string, error) {
	return rm.Render(s, r)
}
