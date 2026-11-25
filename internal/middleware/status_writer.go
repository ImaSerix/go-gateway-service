package middleware

import (
	"bufio"
	"net"
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.ResponseWriter.Write(b)
}

// Hijack hands the connection over for a protocol upgrade (WebSocket).
// ReverseProxy hijacks only after the upstream answered 101 and writes that
// answer itself, so WriteHeader is not called: the status is recorded here.
func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := http.NewResponseController(w.ResponseWriter).Hijack()
	if err == nil {
		w.statusCode = http.StatusSwitchingProtocols
	}

	return conn, rw, err
}

// Unwrap lets http.ResponseController reach the wrapped writer (Flush, deadlines).
func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
