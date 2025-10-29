package auth

import "time"

// ClientConfig is configuration for auth client.
type ClientConfig struct {
	BaseURL      string        // https://myhealthbox.ddns.net
	Path         string        // /api/v1/session
	Timeout      time.Duration // общий timeout запроса (30c)
	CookieName   string        // session_id
	CookieDomain string        // "/"
}
