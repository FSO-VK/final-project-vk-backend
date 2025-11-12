package llmclient

// LLMProvider is an interface for LLMClient.
type LLMProvider interface {
	Query(servicePrompt string) (string, error)
}
