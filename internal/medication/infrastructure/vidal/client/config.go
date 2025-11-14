// Package client is an implementation of client to external service.
package client

import "time"

// Config is a configuration of client.
type Config struct {
	Endpoint string
	APIToken string
	Timeout  time.Duration
}
