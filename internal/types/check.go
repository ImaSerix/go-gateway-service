package types

type CheckName string

var (
	Auth           CheckName = "auth"
	HeaderRequired CheckName = "required_header"
	Inject         CheckName = "inject"
	IPWhiteList    CheckName = "ip_whitelist"
	QueryRequired  CheckName = "required_query"

	//TODO: чеки ниже имеет смысл заменить на middleware
	RateLimitC CheckName = "rate_limit"
	TimeoutC   CheckName = "timeout"
)
