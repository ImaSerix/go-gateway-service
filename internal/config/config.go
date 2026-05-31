package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Root struct {
	Server Server  `yaml:"server"`
	Routes []Route `yaml:"routes"`
}

type Server struct {
	Middlewares []Middleware `yaml:"middlewares"`
}

type Middleware struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

type Route struct {
	Path        string       `yaml:"path"`
	Method      string       `yaml:"method"`
	Middlewares []Middleware `yaml:"middlewares"`
	Checks      []Check      `yaml:"checks"`
	Transforms  Transform    `yaml:"transforms"`
	Upstream    Upstream     `yaml:"upstream"`
}

type Check struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

type Upstream struct {
	Host   string `yaml:"host"`
	Scheme string `yaml:"scheme"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

func LoadConfig(path string) (Root, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return Root{}, fmt.Errorf("loadConfig: %w", err)
	}
	var config Root
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Root{}, fmt.Errorf("loadConfig: %w", err)
	}

	return config, nil
}
