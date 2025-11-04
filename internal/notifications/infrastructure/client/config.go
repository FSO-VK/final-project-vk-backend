package client

import "time"

// PushClient is configuration for notification push client.
type PushClient struct {
	PublicKey  string
	PrivateKey string
	Timeout    time.Duration // общий timeout запроса (30c)
}
