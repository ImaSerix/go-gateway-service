package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// TODO: Сделать нормальны auth middleware, и возможно identity midleware
// TODO: Убрать чеки, которые уже не надо. А также сделать универсальный external policy чек
