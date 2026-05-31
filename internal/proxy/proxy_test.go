package proxy_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/proxy"
)

func TestReverseProxy_ServeHTTP_Success(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if data, _ := io.ReadAll(r.Body); !bytes.Equal(data, []byte("request body")) {
			t.Fatalf("expected request body %s, but got %s", []byte("request body"), data)
		}

		if v := r.Header.Get("X-Username"); v != "usrname" {
			t.Fatalf("expected header %s, but got %s", "usrname", v)
		}

		w.Header().Set("X-Request-ID", "1001")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body"))
	}))
	defer srv.Close()

	target, _ := proxy.NewURL(srv.URL)
	method, _ := proxy.NewMethod("POST")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://nice.url", bytes.NewBuffer([]byte("request body")))

	r.Header.Set("X-Username", "usrname")

	p := proxy.NewReverseProxy(target, method, http.DefaultClient)
	p.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), []byte("response body")) {
		t.Fatalf("expected response body %s, but got %s", []byte("response body"), w.Body.Bytes())
	}

	if v := w.Header().Get("X-Request-ID"); v != "1001" {
		t.Fatalf("expected header %s, but got %s", "1001", v)
	}
}

func TestReverseProxy_ServeHTTP_RequestNil(t *testing.T) {

	srv := httptest.NewServer(nil)
	defer srv.Close()

	target, _ := proxy.NewURL(srv.URL)
	method, _ := proxy.NewMethod("POST")

	w := httptest.NewRecorder()

	p := proxy.NewReverseProxy(target, method, http.DefaultClient)
	p.ServeHTTP(w, nil)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestReverseProxy_ServeHTTP_ClientError(t *testing.T) {

	srv := httptest.NewServer(nil)
	srv.Close()

	target, _ := proxy.NewURL(srv.URL)
	method, _ := proxy.NewMethod("POST")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://nice.url", bytes.NewBuffer([]byte("request body")))

	p := proxy.NewReverseProxy(target, method, http.DefaultClient)
	p.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}
}
