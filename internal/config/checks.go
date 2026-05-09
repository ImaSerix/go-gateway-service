package config

type Store struct {
	Body    map[string]string `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
}

type AuthCheckConfig struct {
	URL            string            `yaml:"url"`
	ForwardHeaders map[string]string `yaml:"forward_headers"`
	Method         string            `yaml:"method"`
	Store          Store             `yaml:"store"`
	ExpectedStatus int               `yaml:"expected_status"`
}

type HeaderRequiredCheckConfig struct {
	Headers []string `yaml:"headers"`
}

type QueryRequiredCheckConfig struct {
	Queries []string `yaml:"queries"`
}

type IPWhiteListCheckConfig struct {
	IP []string `yaml:"ip"`
}

type RateLimitCheckConfig struct {
	Limit  int    `yaml:"limit"`
	Window string `yaml:"window"`
}

type InjectCheckConfig struct {
	Ctx map[string]any `yaml:"ctx"`
}

type TimeoutCheckConfig struct {
	Duration string `yaml:"duration"`
}
