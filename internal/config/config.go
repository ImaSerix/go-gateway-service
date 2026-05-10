package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes []RouteConfig `yaml:"routes"`
}

type RouteConfig struct {
	Path      string          `yaml:"path"`
	Method    string          `yaml:"method"`
	Checks    []CheckConfig   `yaml:"check"`
	Upstream  UpstreamConfig  `yaml:"upstream"`
	Transform TransformConfig `yaml:"transform"`
}

type TransformConfig struct {
	Header map[string]any `yaml:"header"`
	Body   map[string]any `yaml:"body"`
}

type CheckConfig struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

type UpstreamConfig struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loadConfig: %w", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("loadConfig: %w", err)
	}
	return &config, nil
}
