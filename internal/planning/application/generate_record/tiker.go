package generaterecord

import (
	"context"
	"time"
)

// TickerInterface is an interface for ticker.
type TickerInterface interface {
	Run(ctx context.Context, tick func(ctx context.Context))
}

// Ticker implements TickerInterface.
type Ticker struct {
	Interval time.Duration
}

// NewTicker returns a new Ticker.
func NewTicker(interval time.Duration) *Ticker {
	return &Ticker{Interval: interval}
}

// Run runs the ticker.
func (t *Ticker) Run(ctx context.Context, tick func(ctx context.Context)) {
	tk := time.NewTicker(t.Interval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			tick(ctx)
		}
	}
}
