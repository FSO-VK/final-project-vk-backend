package datamatrix

import "time"

// ClientConfig is configuration for DataMatrix client.
type ClientConfig struct {
	BaseURL string        // https://mobile.api.crpt.ru/mobile/check?codeType=datamatrix&code=
	Timeout time.Duration // общий timeout запроса (30c)
}
