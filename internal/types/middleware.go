package types

type MiddlewareName string

var (
	Cors      MiddlewareName = "cors"
	Recovery  MiddlewareName = "recovery"
	RateLimit MiddlewareName = "rate_limit"
	Logging   MiddlewareName = "logging"
	RequestID MiddlewareName = "request_id"
	RealIP    MiddlewareName = "real_ip"
	Timeout   MiddlewareName = "timeout"
	Metric    MiddlewareName = "metric"
	Inject    MiddlewareName = "inject"
)
