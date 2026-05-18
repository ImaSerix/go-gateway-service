package check_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder/check"
	"github.com/ImaSerix/go-gateway-service/internal/config"
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

func TestCheckBuilder_Build(t *testing.T) {

	badError := errors.New("bad error")

	tests := []struct {
		name           string
		registryCalled bool
		registryHas    bool
		factoryCalled  bool
		factoryError   error
		expError       error
	}{
		{
			name:           "success",
			registryCalled: true,
			registryHas:    true,
			factoryCalled:  true,
			expError:       nil,
		},
		{
			name:           "unregistered check",
			registryCalled: true,
			registryHas:    false,
			factoryCalled:  false,
			expError:       check.ErrUnregisteredCheckName,
		},
		{
			name:           "factory error",
			registryCalled: true,
			registryHas:    true,
			factoryCalled:  true,
			factoryError:   badError,
			expError:       badError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			factory := &factoryMock{
				err: test.factoryError,
			}
			reg := &registryMock{
				f:  factory,
				ok: test.registryHas,
			}
			b := check.NewBuilder(reg)

			c, err := b.Build(config.Check{Type: "nice type"})
			if !errors.Is(err, test.expError) {
				t.Fatalf("expected wrapped error %v, but got %v", test.expError, err)
			}
			if reg.called != test.registryCalled {
				t.Fatalf("expected registry called=%v, but called=%v", test.registryCalled, reg.called)
			}
			if factory.called != test.factoryCalled {
				t.Fatalf("expected registry called=%v, but called=%v", test.registryCalled, factory.called)
			}
			if err == nil && c == nil {
				t.Fatal("got no error, but check is nil")
			}
			if err != nil && c != nil {
				t.Fatal("got error, but check is not nil")
			}
		})
	}
}

func TestCheckBuilder_BuildMany_Success(t *testing.T) {

	cfgs := []config.Check{
		config.Check{Type: "check1"},
		config.Check{Type: "check2"},
	}

	factory := &factoryMock{}
	reg := &registryMock{
		f:  factory,
		ok: true,
	}

	b := check.NewBuilder(reg)
	c, err := b.BuildMany(cfgs)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}
	if len(c) != 2 {
		t.Fatalf("expected 2 checks, but got %d", len(c))
	}

}

func TestCheckBuilder_BuildMany_Error(t *testing.T) {

	cfgs := []config.Check{
		config.Check{Type: "bad check"},
		config.Check{Type: "check2"},
	}

	factory := &factoryMock{}
	reg := &registryMock{
		f:  factory,
		ok: true,
	}

	b := check.NewBuilder(reg)
	c, err := b.BuildMany(cfgs)
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
	if len(c) != 0 {
		t.Fatalf("expected empty checks, but got %d", len(c))
	}

}
