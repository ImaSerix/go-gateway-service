package config

type Store map[string]string

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
