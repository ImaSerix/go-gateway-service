package endpoint

import (
	"io"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type Endpoint struct {
	Path         Path
	Method       Method
	Target       URL
	TargetMethod Method
}

func NewEndpoint(path Path, method Method, target URL, targetMethod Method) *Endpoint {
	return &Endpoint{
		Path:         path,
		Method:       method,
		Target:       target,
		TargetMethod: targetMethod,
	}
}

func NewEndpointFromConfig(cfg *config.RouteConfig) (*Endpoint, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	path, err := NewPath(cfg.Path)
	if err != nil {
		return nil, err
	}
	method, err := NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}
	target, err := NewURL(cfg.Upstream.URL)
	if err != nil {
		return nil, err
	}

	tm := cfg.Upstream.Method
	if tm == "" {
		tm = cfg.Method
	}

	targetMethod, err := NewMethod(tm)
	if err != nil {
		return nil, err
	}
	return &Endpoint{
		Path:         path,
		Method:       method,
		Target:       target,
		TargetMethod: targetMethod,
	}, nil
}

func (e *Endpoint) matchMethod(method string) bool {
	return e.Method == Method(method)
}

func (e *Endpoint) Pattern() string {
	return string(e.Method) + " " + string(e.Path)
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !e.matchMethod(r.Method) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var resp *http.Response
	var err error
	switch e.TargetMethod {
	case GET:
		resp, err = http.Get(string(e.Target))
		defer resp.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	case POST:
		// TODO: На данный момент тип контента один
		resp, err = http.Post(string(e.Target), "text/plain", nil)
		defer resp.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
}
