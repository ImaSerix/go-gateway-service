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
	Middleware []Middleware `yaml:"middleware"`
}

type Middleware struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

type Route struct {
	Path       string       `yaml:"path"`
	Method     string       `yaml:"method"`
	Middleware []Middleware `yaml:"middleware"`
	Checks     []Check      `yaml:"checks"`
	Upstream   Upstream     `yaml:"upstream"`
	Transform  Transform    `yaml:"transform"`
}

type Transform struct {
	Header map[string]string `yaml:"header"`
	Body   map[string]any    `yaml:"body"`
}

type Check struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

type Upstream struct {
	URL    string `yaml:"url"`
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
