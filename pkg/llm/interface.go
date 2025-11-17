package llm

// Provider is an interface for LLMClient.
type Provider interface {
	Query(servicePrompt string) (string, error)
	UseSystemPrompt(servicePrompt string, additionalPromptPath string) (string, error)
}
