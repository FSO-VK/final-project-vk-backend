package daemon

import "time"

// ClientConfig is configuration for generate record daemon.
type ClientConfig struct {
	// createdShift is the offset from 00:00 when records are generated.
	// At 00:00 + createdShift, all records for that day are created. (basically 24h - today creating for the next day)
	CreationShift  time.Duration `koanf:"creation_shift"`
	TickerInterval time.Duration `koanf:"ticker_interval"`
}
