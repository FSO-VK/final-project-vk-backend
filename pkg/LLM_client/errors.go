package llmclient

import "errors"

var (
	// ErrAuthFailed is returned when the LLM server returns an unsuccessful auth status.
	ErrAuthFailed = errors.New("llmclient: auth failed")
	// ErrEmptyAccessToken is returned when the access token is missing in the response.
	ErrEmptyAccessToken = errors.New("llmclient: empty access token in response")
	// ErrEmptyPrompt is returned when the API request was called with an empty prompt.
	ErrEmptyPrompt = errors.New("llmclient: empty prompt")
	// ErrAPIRequestFailed is returned when the LLM API responds with non-200 OK status.
	ErrAPIRequestFailed = errors.New("llmclient: API request failed")
	// ErrAPIError is returned when the API returns an error object ("error" field in JSON).
	ErrAPIError = errors.New("llmclient: API returned error")
	// ErrEmptyResponse is returned when the response body is empty or contains no data.
	ErrEmptyResponse = errors.New("llmclient: empty response")
)
