package config

type Store struct {
	Headers    map[string]string `yaml:"headers"`
	BodyFields map[string]string `yaml:"body_fields"`
}

// type AuthCheck struct {
// 	URL            string            `yaml:"url"`
// 	ForwardHeaders map[string]string `yaml:"forward_headers"`
// 	Method         string            `yaml:"method"`
// 	Store          Store             `yaml:"store"`
// 	ExpectedStatus int               `yaml:"expected_status"`
// }

type PolicyCheck struct {
	Transform      Transform `yaml:"transform"`
	Upstream       Upstream  `yaml:"upstream"`
	ExpectedStatus int       `yaml:"expected_status"`
	Store          Store     `yaml:"store"`
}

type HeaderRequiredCheck struct {
	Headers []string `yaml:"headers"`
}

type QueryRequiredCheck struct {
	QueryParams []string `yaml:"query_params"`
}

type IPWhiteListCheck struct {
	IPs []string `yaml:"ips"`
}

// type RateLimitCheck struct {
// 	Limit  int    `yaml:"limit"`
// 	Window string `yaml:"window"`
// }

// type InjectCheck struct {
// 	Ctx map[string]any `yaml:"ctx"`
// }

// type TimeoutCheck struct {
// 	Duration string `yaml:"duration"`
// }
