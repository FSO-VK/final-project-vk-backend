package llmclient

import "time"

// ClientConfig is configuration for DataMatrix client.
type ClientConfig struct {
	LLMBaseURL   string // "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	AuthBaseURL  string // "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	ClientID     string
	ClientSecret string
	ModelName    string        // "GigaChat"
	Role         string        // "user"
	Temperature  float64       // 0.1
	MaxTokens    int           // 2000
	Timeout      time.Duration // общий timeout запроса (30c)
}
