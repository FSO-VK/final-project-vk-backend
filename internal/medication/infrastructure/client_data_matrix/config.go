package clientdatamatrix

import "time"

// ClientConfig is configuration for DataMatrix client.
type ClientConfig struct {
	Host    string        // https://mobile.api.crpt.ru/mobile/check?codeType=datamatrix&code=
	Timeout time.Duration // общий timeout запроса (30c)
}
