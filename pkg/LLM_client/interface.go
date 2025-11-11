package llmclient

// LLMClientProvider is an interface for LLMClient.
type LLMClientProvider interface {
	Query(servicePrompt string) (string, error)
}
