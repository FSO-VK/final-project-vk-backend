package daemon

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// DaemonProvider is an interface for ticker.
type DaemonProvider interface {
	Run(ctx context.Context)
	Register(f func(ctx context.Context) error)
}

// Daemon implements DaemonProvider.
type Daemon struct {
	TickerInterval time.Duration
	StartAt        *time.Time // nil - start immediately
	logger         *logrus.Entry
}

// NewDaemon returns a new Daemon.
func NewDaemon(
	interval time.Duration,
	startAt *time.Time,
	logger *logrus.Entry,
) *Daemon {
	return &Daemon{
		TickerInterval: interval,
		StartAt:        startAt,
		logger:         logger,
	}
}

// RunTicker runs the ticker.
func (t *Daemon) RunTicker(ctx context.Context, tick func(ctx context.Context) error) {
	if t.StartAt != nil {
		wait := time.Until(*t.StartAt)
		if wait > 0 {
			t.logger.Infof("wait for %s", wait)
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
		}
		t.logger.Info("execute func")
		if err := tick(ctx); err != nil {
			t.logger.Error("err on func", err)
		}
	}

	tk := time.NewTicker(t.TickerInterval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			t.logger.Info("execute func")
			if err := tick(ctx); err != nil {
				t.logger.Error("err on func", err)
			}
		}
	}
}

// Run runs the daemon.
func (d *Daemon) Run(ctx context.Context, f func(ctx context.Context) error) {
	d.RunTicker(ctx, f)
}
