package proxy

// REPLACED BY httputils.SingleHostReverseProxy

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
)

type ReverseProxy struct {
	target URL
	method Method
	client *http.Client
}

func NewReverseProxy(target URL, method Method, client *http.Client) *ReverseProxy {
	return &ReverseProxy{
		target: target,
		method: method,
		client: client,
	}
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r == nil {
		slog.Log(context.Background(), slog.LevelError, "nil request")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(string(p.method), string(p.target), nil)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "error on new request", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	for k, v := range r.Header {
		req.Header[k] = v
	}

	b, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(b))
	req.Body = io.NopCloser(bytes.NewBuffer(b))
	req.ContentLength = r.ContentLength

	resp, err := p.client.Do(req)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "error on request", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, b := range resp.Header {
		w.Header()[k] = b
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
