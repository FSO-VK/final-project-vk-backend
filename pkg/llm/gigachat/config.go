package gigachat

import "time"

// ClientConfig is configuration for llm provider client.
type ClientConfig struct {
	LLMBaseURL           string        `koanf:"llm_base_url"`
	AuthBaseURL          string        `koanf:"auth_base_url"`
	ClientID             string        `koanf:"client_id"`
	ClientSecret         string        `koanf:"client_secret"`
	ModelName            string        `koanf:"model_name"`
	Role                 string        `koanf:"role"`
	Temperature          float64       `koanf:"temperature"`
	MaxTokens            int           `koanf:"max_tokens"`
	Timeout              time.Duration `koanf:"timeout"`
	BaseSystemPromptPath string        `koanf:"base_system_prompt_path"`
}
