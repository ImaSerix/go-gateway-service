package middleware_test

import (
	"errors"
	"testing"

	"github.com/ImaSerix/go-gateway-service/internal/builder/middleware"
	"github.com/ImaSerix/go-gateway-service/internal/config"
)

func TestMiddlewareBuilder_Build(t *testing.T) {

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
			name:           "unregistered middleware",
			registryCalled: true,
			registryHas:    false,
			factoryCalled:  false,
			expError:       middleware.ErrUnregisteredMiddlewareType,
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

			b := middleware.NewBuilder(reg)

			c, err := b.Build(config.Middleware{Type: "nice type"})
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
				t.Fatal("got no error, but middlewar is nil")
			}
			if err != nil && c != nil {
				t.Fatal("got error, but middleware is not nil")
			}
		})
	}
}

func TestMiddlewareBuilder_BuildMany_Success(t *testing.T) {

	cfgs := []config.Middleware{
		config.Middleware{Type: "middleaware1"},
		config.Middleware{Type: "middleaware2"},
	}

	factory := &factoryMock{}
	reg := &registryMock{
		f:  factory,
		ok: true,
	}

	b := middleware.NewBuilder(reg)
	c, err := b.BuildMany(cfgs)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}
	if len(c) != 2 {
		t.Fatalf("expected 2 checks, but got %d", len(c))
	}

}

func TestMiddlewareBuilder_BuildMany_Error(t *testing.T) {

	cfgs := []config.Middleware{
		config.Middleware{Type: "bad middleware"},
		config.Middleware{Type: "middleware2"},
	}

	factory := &factoryMock{}
	reg := &registryMock{
		f:  factory,
		ok: true,
	}

	b := middleware.NewBuilder(reg)
	c, err := b.BuildMany(cfgs)
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
	if len(c) != 0 {
		t.Fatalf("expected empty checks, but got %d", len(c))
	}

}
