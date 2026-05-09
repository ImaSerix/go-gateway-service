package endpoint

import (
	"io"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

type Endpoint struct {
	Path     Path
	Method   Method
	Checks   []Check
	Upstream *Upstream
}

func NewEndpoint(path Path, method Method, target URL, targetMethod Method) *Endpoint {
	return &Endpoint{
		Path:   path,
		Method: method,
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

	if cfg.Upstream.Method == "" {
		cfg.Upstream.Method = cfg.Method
	}

	upstream, err := NewUpstreamFromConfig(&cfg.Upstream)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		Path:     path,
		Method:   method,
		Upstream: upstream,
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

	//TODO: нужно сделать с контекстом, иначе херня
	var resp *http.Response
	var err error
	switch e.Upstream.Method {
	case GET:
		resp, err = http.Get(string(e.Upstream.URL))
		defer resp.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	case POST:
		// TODO: На данный момент тип контента один
		resp, err = http.Post(string(e.Upstream.URL), "text/plain", nil)
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
