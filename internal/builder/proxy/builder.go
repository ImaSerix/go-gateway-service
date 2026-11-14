package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/httperror"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/proxy"
	"github.com/ImaSerix/go-gateway-service/internal/renderer"
)

// rewriteErrKey is the context key of the error that Rewrite hands over to ErrorHandler.
type rewriteErrKey struct{}

type ReverseBuilder struct {
	client *http.Client
	render renderer.Renderer
}

func NewBuilder(client *http.Client, render renderer.Renderer) *ReverseBuilder {
	return &ReverseBuilder{
		client: client,
		render: render,
	}
}

func (rb *ReverseBuilder) Build(cfg config.Upstream) (pipeline.Proxy, error) {

	u, err := proxy.MakeURL(cfg.Scheme, cfg.Host, cfg.Path)
	if err != nil {
		return nil, err
	}

	m, err := proxy.NewMethod(cfg.Method)
	if err != nil {
		return nil, err
	}

	var transport http.RoundTripper
	if rb.client != nil {
		transport = rb.client.Transport
	}

	p := &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {

			renderedURL, err := rb.render.Render(rawTarget(u), pr.In)

			var target *url.URL
			if err == nil {
				target, err = url.Parse(renderedURL)
			}
			if err != nil {
				// Rewrite cannot return an error: the outgoing request is left without
				// an address and carries the error, ErrorHandler answers with 500.
				pr.Out = pr.Out.WithContext(context.WithValue(pr.Out.Context(), rewriteErrKey{}, err))
				pr.Out.URL.Scheme = ""
				pr.Out.URL.Host = ""
				return
			}

			pr.SetURL(target)

			// SetURL joins the target path with the incoming one,
			// a configured path replaces the incoming path instead.
			if target.Path != "" {
				pr.Out.URL.Path = target.Path
				pr.Out.URL.RawPath = target.RawPath
			}

			pr.Out.Method = string(m)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if rewriteErr, ok := r.Context().Value(rewriteErrKey{}).(error); ok {
				slog.Error("upstream url build failed", "method", r.Method, "path", r.URL.Path, "err", rewriteErr)
				httperror.Write(w, http.StatusInternalServerError, httperror.CodeInternalError, "internal error")
				return
			}

			slog.Error("upstream request failed", "method", r.Method, "path", r.URL.Path, "err", err)
			httperror.Write(w, http.StatusBadGateway, httperror.CodeUpstreamUnavailable, "upstream service unavailable")
		},
	}

	return p, nil
}

func rawTarget(u *url.URL) string {
	raw := u.Scheme + "://" + u.Host + u.Path
	if u.RawQuery != "" {
		raw += "?" + u.RawQuery
	}

	return raw
}
