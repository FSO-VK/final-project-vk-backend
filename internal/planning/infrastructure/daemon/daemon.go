package daemon

import (
	"context"
	"time"
)

// DaemonProvider is an interface for ticker.
type DaemonProvider interface {
	Run(ctx context.Context)
	Register(f func(ctx context.Context) error)
}

// Daemon implements DaemonProvider.
type Daemon struct {
	TickerInterval time.Duration
}

// NewDaemon returns a new Daemon.
func NewDaemon(interval time.Duration) *Daemon {
	return &Daemon{
		TickerInterval: interval,
	}
}

// RunTicker runs the ticker.
func (t *Daemon) RunTicker(ctx context.Context, tick func(ctx context.Context) error) {
	tk := time.NewTicker(t.TickerInterval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			_ = tick(ctx)
		}
	}
}

// Run runs the daemon.
func (d *Daemon) Run(ctx context.Context, f func(ctx context.Context) error) {
	d.RunTicker(ctx, f)
}
