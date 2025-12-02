package notificationclient

import "time"

// ClientConfig is configuration for medication client.
type ClientConfig struct {
	Host    string
	Method  string
	Timeout time.Duration // общий timeout запроса (30c)
}
