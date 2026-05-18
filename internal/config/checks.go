package config

type Store struct {
	Body    map[string]string `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
}

type AuthCheck struct {
	URL            string            `yaml:"url"`
	ForwardHeaders map[string]string `yaml:"forward_headers"`
	Method         string            `yaml:"method"`
	Store          Store             `yaml:"store"`
	ExpectedStatus int               `yaml:"expected_status"`
}

type HeaderRequiredCheck struct {
	Header []string `yaml:"header"`
}

type QueryRequiredCheck struct {
	Query []string `yaml:"query"`
}

type IPWhiteListCheck struct {
	IP []string `yaml:"ip"`
}

type RateLimitCheck struct {
	Limit  int    `yaml:"limit"`
	Window string `yaml:"window"`
}

type InjectCheck struct {
	Ctx map[string]any `yaml:"ctx"`
}

type TimeoutCheck struct {
	Duration string `yaml:"duration"`
}
