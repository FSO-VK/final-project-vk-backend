package llmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// LLMClient implements LLMClientProvider.
type LLMClient struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

// NewLLMClient creates a new LLMClientProvider.
func NewLLMClient(cfg ClientConfig, logger *logrus.Entry) *LLMClient {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &LLMClient{client: client, cfg: cfg, logger: logger}
}

// Query implements LLMClientProvider interface.
func (h *LLMClient) Query(servicePrompt string) (string, error) {
	if servicePrompt == "" {
		return "", ErrEmptyPrompt
	}
	token, err := getGigaChatToken(h.cfg.ClientID,
		h.cfg.ClientSecret,
		h.cfg.AuthBaseURL,
		h.cfg.Timeout,
	)
	if err != nil {
		h.logger.Errorf("Failed to get token: %v", err)
		return "", err
	}

	template, err := template.ParseFiles("templates/prompt.tmpl")
	if err != nil {
		h.logger.Errorf("Failed to parse template: %v", err)
		return "", err
	}
	data := PromptData{Prompt: servicePrompt}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		h.logger.Errorf("Failed to execute template: %v", err)
		return "", err
	}

	fullPrompt := buf.String()

	response, err := askGigaChat(h.cfg, token, fullPrompt)
	if err != nil {
		h.logger.Errorf("Failed to get response: %v", err)
		return "", err
	}
	h.logger.Info("Correct response from GigaChat")

	if err = validateIfJSON(response); err != nil {
		h.logger.Errorf("not json in response: %v", err)
		return "", err
	}
	return "", nil
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
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		cfg.LLMBaseURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := createHTTPClient(cfg.Timeout)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	response, err := parseGigaChatResponse(resp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	return response, nil
}

func parseGigaChatResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrAPIError
	}

	var result GigaChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error.Message != "" {
		return "", ErrAPIRequestFailed
	}

	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		return result.Choices[0].Message.Content, nil
	}

	return "", ErrEmptyResponse
}

func validateIfJSON(response string) error {
	resp := strings.TrimSpace(response)
	if resp == "" {
		return ErrEmptyResponse
	}

	var js map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &js); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}
