package check

import (
	"fmt"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
)

type Builder struct {
	registry Registry
}

func NewBuilder(registry Registry) *Builder {
	return &Builder{
		registry: registry,
	}
}

func (b *Builder) Build(cfg config.Check) (pipeline.Checker, error) {

	// Используется регистер для более удобного и юнифицированного и понятного создания check
	f, ok := b.registry.Get(types.CheckName(cfg.Type))
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnregisteredCheckName, cfg.Type)
	}

	m, err := f.Create(cfg.Config)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (b *Builder) BuildMany(cfg []config.Check) ([]pipeline.Checker, error) {
	result := make([]pipeline.Checker, 0, len(cfg))

	for _, c := range cfg {
		check, err := b.Build(c)
		if err != nil {
			return nil, err
		}

		result = append(result, check)
	}

	return result, nil
}
