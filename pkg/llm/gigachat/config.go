package gigachat

import "time"

// ClientConfig is configuration for llm provider client.
type ClientConfig struct {
	LLMBaseURL           string
	AuthBaseURL          string
	ClientID             string
	ClientSecret         string
	ModelName            string
	Role                 string
	Temperature          float64
	MaxTokens            int
	Timeout              time.Duration
	BaseSystemPromptPath string
}
