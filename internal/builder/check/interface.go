package check

import (
	"github.com/ImaSerix/go-gateway-service/internal/client"
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

type TransformBuilder interface {
	BuildMany(config.Transform) ([]pipeline.Transformer, error)
}

type StoreBuilder interface {
	Build(config.Store) (pipeline.Store, error)
}

type ClientBuilder interface {
	Build(config.Upstream) (*client.Upstream, error)
}

type Factory interface {
	Create(raw yaml.Node) (pipeline.Checker, error)
}

type Registry interface {
	Get(key types.CheckName) (Factory, bool)
}
