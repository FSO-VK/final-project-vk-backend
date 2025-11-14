package client

import "time"

type Config struct {
	Endpoint string
	APIToken string
	Timeout  time.Duration
}
