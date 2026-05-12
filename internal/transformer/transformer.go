package transformer

import (
	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

// TODO: Сделать фабрику трансформов, а также реализовать в целом разные трансформы для хэдэра, боди

func TransformersFactory(cfg *config.TransformConfig) ([]pipeline.Transformer, error) {

	res := []pipeline.Transformer{}

	if cfg.Body != nil {
		res = append(res, NewBodyTransformer(cfg.Body))
	}

	if cfg.Header != nil {
		res = append(res, NewHeaderTransformer(cfg.Header))
	}

	return res, nil
}
