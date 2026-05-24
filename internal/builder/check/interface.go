package check

import (
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

type Factory interface {
	Create(raw yaml.Node) (pipeline.Checker, error)
}

type Registry interface {
	Get(key types.CheckName) (Factory, bool)
}
