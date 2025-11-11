package llmclient

// PromptData is the data for the prompt template.
type PromptData struct {
	Prompt string
}

// GigaChatRequest is the request body for the GigaChat API.
type GigaChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	//nolint:tagliatelle
	MaxTokens int `json:"max_tokens,omitempty"` // да у гигачата json с нижними подчеркиваниями а линтер ругается
}

// Message is a message in the GigaChat API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GigaChatResponse is the response body for the GigaChat API.
type GigaChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type GigaChatAuth struct {
	//nolint:tagliatelle
	AccessToken string `json:"access_token"` // да у гигачата json с нижними подчеркиваниями а линтер ругается
	//nolint:tagliatelle
	ExpiresAt int64 `json:"expires_at"`
}
