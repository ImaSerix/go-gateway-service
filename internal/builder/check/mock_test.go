package check_test

import (
	"context"
	"net/http"

	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/pipeline"
	"github.com/ImaSerix/go-gateway-service/internal/types"
	"gopkg.in/yaml.v3"
)

type checkMock struct{}

func (cm checkMock) Execute(ctx context.Context, r *http.Request) (context.Context, error) {
	return context.Background(), nil
}

type factoryMock struct {
	called bool
	err    error
}

func (fm *factoryMock) Create(raw yaml.Node) (pipeline.Checker, error) {
	fm.called = true
	if fm.err != nil {
		return nil, fm.err
	}
	return &checkMock{}, fm.err
}

type registryMock struct {
	f      *factoryMock
	ok     bool
	called bool
}

func (rm *registryMock) Get(key types.CheckName) (check.Factory, bool) {
	rm.called = true
	//TODO: Если честно херня полная, но мб и так норм, лол
	if key == "bad check" {
		return nil, false
	}
	return rm.f, rm.ok
}

type rendererMock struct {
	res string
	err error
}

func (rm *rendererMock) Render(s string, r *http.Request) (string, error) {
	return rm.res, rm.err
}
