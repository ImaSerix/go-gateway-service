package checks

import (
	"context"
	"net/http"
	"time"

	"github.com/ImaSerix/go-gateway-service/internal/config"
)

const TimeoutCheckType = "timeout"

type TimeoutCheck struct {
	duration time.Duration
}

func NewTimeoutCheck(cfg *config.TimeoutCheckConfig) (*TimeoutCheck, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	d, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, ErrInvalidDuration
	}

	return &TimeoutCheck{
		duration: d,
	}, nil
}

func (c *TimeoutCheck) Execute(ctx context.Context, r *http.Request) (context.Context, error) {

	if r == nil {
		return ctx, ErrNilRequest
	}

	//TODO: Сейчас cancel игнорируется может быть позже имеет смысл перенести в middlewarе
	ctx, _ = context.WithTimeout(ctx, c.duration)
	return ctx, nil
}
