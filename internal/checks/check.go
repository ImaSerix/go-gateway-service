package checks

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/config"
	"github.com/ImaSerix/go-gateway-service/internal/endpoint"
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

func CheckFactory(cfg *config.Check, client *http.Client) (endpoint.Check, error) {
	switch cfg.Type {
	case AuthCheckType:
		var typedCfg config.AuthCheckConfig
		err := cfg.Config.Decode(&typedCfg)
		if err != nil {
			return nil, ErrInvalidConfig
		}
		c, err := NewAuthCheck(&typedCfg, client)
		return c, err
	default:
		return nil, ErrUnsupportedType
	}
}
