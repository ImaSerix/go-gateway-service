// Package httperror writes the errors that the gateway itself returns to the client.
// Responses of upstream services are not touched.
package httperror

import (
	"encoding/json"
	"net/http"
)

const (
	CodeBadRequest          = "bad_request"
	CodeUnauthorized        = "unauthorized"
	CodeAuthUnavailable     = "auth_unavailable"
	CodeNotFound            = "not_found"
	CodeMethodNotAllowed    = "method_not_allowed"
	CodeForbidden           = "forbidden"
	CodeTooManyRequests     = "too_many_requests"
	CodeUpstreamUnavailable = "upstream_unavailable"
	CodeInternalError       = "internal_error"
)

type body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Write sends {"code": "...", "message": "..."} with the given status.
// The message must not contain tokens, internal addresses or internal error text.
func Write(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body{Code: code, Message: message})
}
