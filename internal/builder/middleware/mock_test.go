package middleware_test

import (
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

func middlewareMock(next http.Handler) http.Handler {
	return next
}

type factoryMock struct {
	called bool
	err    error
}

func (fm *factoryMock) Create(raw yaml.Node) (pipeline.Middleware, error) {
	fm.called = true
	if fm.err != nil {
		return nil, fm.err
	}
	return middlewareMock, fm.err
}

type registryMock struct {
	f      *factoryMock
	ok     bool
	called bool
}

func (rm *registryMock) Get(key types.MiddlewareName) (middleware.Factory, bool) {
	rm.called = true
	//TODO: Если честно херня полная, но мб и так норм, лол
	if key == "bad middleware" {
		return nil, false
	}
	return rm.f, rm.ok
}
