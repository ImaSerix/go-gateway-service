package endpoint

import "net/http"

type Check interface {
	Execute(r *http.Request) (bool, error)
}
