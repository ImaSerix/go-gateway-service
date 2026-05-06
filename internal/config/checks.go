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
