package check

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/ImaSerix/go-gateway-service/internal/config"
// )

// type Timeout struct {
// 	duration time.Duration
// }

// func NewTimeout(cfg config.TimeoutCheck) (*Timeout, error) {

// 	d, err := time.ParseDuration(cfg.Duration)
// 	if err != nil {
// 		return nil, ErrInvalidDuration
// 	}

// 	return &Timeout{
// 		duration: d,
// 	}, nil
// }

// func (c *Timeout) Execute(ctx context.Context, r *http.Request) (context.Context, error) {

// 	if r == nil {
// 		return ctx, ErrNilRequest
// 	}

// 	//TODO: Сейчас cancel игнорируется может быть позже имеет смысл перенести в middlewarе и там уже делать что-то с этим
// 	ctx, _ = context.WithTimeout(ctx, c.duration)
// 	return ctx, nil
// }
