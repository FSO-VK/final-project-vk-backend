package datamatrix

import "time"

// ClientConfig is configuration for DataMatrix client.
type ClientConfig struct {
	Host    string
	Timeout time.Duration // общий timeout запроса (30c)
}
