package transformer

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

// TODO: Сделать фабрику трансформов, а также реализовать в целом разные трансформы для хэдэра, боди
func TransformerFactory(cfg *config.TransformConfig) ([]pipeline.Transformer, error) {

}
