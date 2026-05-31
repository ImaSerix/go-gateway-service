package config

import "gopkg.in/yaml.v3"

type Transform map[string]yaml.Node

type QueryParamsTransform map[string]any
type HeadersTransform map[string]string
type BodyFieldsTransform map[string]any
