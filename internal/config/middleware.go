package config

// --- CORS Middleware ---
type CORSMiddleaware_Allowed struct {
	Origin []string `yaml:"origin"`
	Method []string `yaml:"method"`
	Header []string `yaml:"header"`
}

type CORSMiddleaware struct {
	Allowed CORSMiddleaware_Allowed `yaml:"allowed"`
}

// Структура выглядит примерно так:

// - type: cors
//   config:
// 	allowed:
// 		origin:
// 			- http://allowed.origin
// 			- http://allowed.origin2
// 		method:
// 			- GET
// 			- POST
// 			- PUT
// 		Header:
// 			- Content-Type
// 			- Authorization

// ---

// --- Rate limit middleware ---
type RateLimitMiddleware struct {
	Limit  int    `yaml:"limit"`
	Window string `yaml:"window"`
}

// Структура выглядит примерно так:

// - type: cors
//   config:
// 		limit: 20
//		window: 1m

// ---

type TimeoutMiddleware struct {
	Duration string `yaml:"duration"`
}

// Структура выглядит примерно так:

// - type: timeout
//   config:
// 		duration: 2s

// ---
