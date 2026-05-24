package transformer_test

import "net/http"

type rendererMock struct {
	res string
	err error
}

func (rm *rendererMock) Render(s string, r *http.Request) (string, error) {
	return rm.res, rm.err
}
