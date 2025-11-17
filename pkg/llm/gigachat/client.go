// Package gigachat provides a client for the GigaChat API.
package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
)

var (
	// ErrEmptyPrompt is returned when the API request was called with an empty prompt.
	ErrEmptyPrompt = errors.New("empty prompt")
	// ErrAPIRequestFailed is returned when the LLM API responds with non-200 OK status.
	ErrAPIRequestFailed = errors.New("API request failed")
	// ErrAPIError is returned when the API returns an error object ("error" field in JSON).
	ErrAPIError = errors.New("API returned error")
	// ErrEmptyResponse is returned when the response body is empty or contains no data.
	ErrEmptyResponse = errors.New("empty response")
	// ErrFailedToGetToken is returned when the access token is missing in the response.
	ErrFailedToGetToken = errors.New("failed to get token")
	// ErrWithSystemPrompt is returned when the access token is missing in the response.
	ErrWithSystemPrompt = errors.New("failed to get token")
	// ErrInvalidResponse is returned when the response is invalid.
	ErrInvalidResponse = errors.New("invalid response")
)

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
	MaxTokens int `json:"max_tokens,omitempty"`
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

// GigachatLLMProvider implements LLMClientProvider.
type GigachatLLMProvider struct {
	client *http.Client
	cfg    ClientConfig
}

// NewGigachatLLMProvider creates a new LLMClientProvider.
func NewGigachatLLMProvider(cfg ClientConfig) *GigachatLLMProvider {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &GigachatLLMProvider{client: client, cfg: cfg}
}

// Query implements LLMClientProvider interface.
func (h *GigachatLLMProvider) Query(servicePrompt string) (string, error) {
	if servicePrompt == "" {
		return "", ErrEmptyPrompt
	}
	token, err := getGigaChatToken(h.cfg.ClientID,
		h.cfg.ClientSecret,
		h.cfg.AuthBaseURL,
		h.cfg.Timeout,
	)
	if err != nil {
		return "", ErrFailedToGetToken
	}

	template, err := template.ParseFiles("./gigachat_prompt/prompt.tmpl")
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	data := PromptData{Prompt: servicePrompt}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	fullPrompt := buf.String()

	response, err := askGigaChat(h.cfg, token, fullPrompt)
	if err != nil {
		return "", ErrInvalidResponse
	}

	if response != "" {
		return "", ErrInvalidResponse
	}
	return response, nil
}

func askGigaChat(cfg ClientConfig, token string, prompt string) (string, error) {
	request := GigaChatRequest{
		Model: cfg.ModelName,
		Messages: []Message{
			{Role: cfg.Role, Content: prompt},
		},
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", ErrInvalidResponse
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		cfg.LLMBaseURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", ErrAPIRequestFailed
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := createHTTPClient(cfg.Timeout)

	resp, err := client.Do(req)
	if err != nil {
		return "", ErrAPIRequestFailed
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	response, err := parseGigaChatResponse(resp)
	if err != nil {
		return "", ErrInvalidResponse
	}
	return response, nil
}

func parseGigaChatResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrAPIRequestFailed
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrAPIError
	}

	var result GigaChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", ErrAPIRequestFailed
	}

	if result.Error.Message != "" {
		return "", ErrAPIRequestFailed
	}

	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		return result.Choices[0].Message.Content, nil
	}

	return "", ErrEmptyResponse
}

// UseSystemPrompt is Query but with additional system prompt.
func (h *GigachatLLMProvider) UseSystemPrompt(
	servicePrompt string,
	additionalPromptPath string,
) (string, error) {
	if servicePrompt == "" {
		return "", ErrEmptyPrompt
	}

	template, err := template.ParseFiles(additionalPromptPath)
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	data := PromptData{Prompt: servicePrompt}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	fullPrompt := buf.String()

	response, err := h.Query(fullPrompt)
	if err != nil {
		return "", ErrInvalidResponse
	}

	return response, nil
}
