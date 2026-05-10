package checks

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
)

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	INVALID Method = "INVALID"
)

func NewMethod(m string) (Method, error) {
	switch m {
	case string(GET):
		return GET, nil
	case string(POST):
		return POST, nil
	default:
		return INVALID, ErrInvalidMethod
	}
}

func ChecksFactory(cfg []config.CheckConfig, client *http.Client) ([]pipeline.Checker, error) {

	res := []pipeline.Checker{}
	for _, c := range cfg {
		checker, err := CheckFactory(c, client)
		if err != nil {
			return nil, err
		}
		res = append(res, checker)
	}
	return res, nil

}

func CheckFactory(cfg config.CheckConfig, client *http.Client) (pipeline.Checker, error) {

	switch cfg.Type {
	case AuthCheckType:

		var typedCfg config.AuthCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewAuthCheck(&typedCfg, client)
		return c, err

	case HeaderRequiredCheckType:

		var typedCfg config.HeaderRequiredCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewHeaderRequiredCheck(&typedCfg)
		return c, err

	case QueryRequiredCheckType:

		var typedCfg config.QueryRequiredCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewQueryRequiredCheck(&typedCfg)
		return c, err

	case IPWhiteListCheckType:

		var typedCfg config.IPWhiteListCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewIPWhiteListCheck(&typedCfg)
		return c, err

	case RateLimitCheckType:

		var typedCfg config.RateLimitCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewRateLimitCheck(&typedCfg)
		return c, err

	case InjectCheckType:

		var typedCfg config.InjectCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewInjectCheck(&typedCfg)
		return c, err

	case TimeoutCheckType:

		var typedCfg config.TimeoutCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewTimeoutCheck(&typedCfg)
		return c, err

	default:
		return nil, ErrUnsupportedType
	}

}
