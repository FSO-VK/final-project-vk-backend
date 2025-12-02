package daemon

import (
	"context"
	"time"
)

// DaemonInterface is an interface for ticker.
type DaemonInterface interface {
	Run(ctx context.Context)
	Register(f func(ctx context.Context) error)
}

// Daemon implements TickerInterface.
type Daemon struct {
	TickerInterval time.Duration
	Function       func(ctx context.Context) error
}

// NewDaemon returns a new Ticker.
func NewDaemon(interval time.Duration) *Daemon {
	return &Daemon{
		TickerInterval: interval,
	}
}

// RunTicker runs the ticker.
func (t *Daemon) RunTicker(ctx context.Context, tick func(ctx context.Context)) {
	tk := time.NewTicker(t.TickerInterval)
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

// Run runs the daemon.
func (d *Daemon) Run(ctx context.Context) {
	d.RunTicker(ctx, func(ctx context.Context) {
		_ = d.Function(ctx)
	})
}

func (d *Daemon) Register(f func(ctx context.Context) error) {
	d.Function = f
}
