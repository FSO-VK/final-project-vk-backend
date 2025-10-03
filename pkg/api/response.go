// Package api provides common interaction protocol between frontend and backend
package api

import "net/http"

// Response is a generic API response wrapper that encapsulates the HTTP status code,
// a response body of any type, and an optional error field.
// T represents the type of the response body.
// The Error field is omitted from the JSON output if it is empty.
type Response[T any] struct {
	StatusCode int       `json:"statusCode"`
	Body       T         `json:"body"`
	Error      ErrorType `json:"error,omitempty"`
}

// NewResponse creates and returns a new Response instance with the specified status code, body, and error message.
// T is a generic type parameter representing the type of the response body.
// Parameters:
//   - statusCode: HTTP status code to be set in the response.
//   - body: The response body of generic type T.
//   - errorMsg: An error message of type ErrorType to be included in the response.
//
// Returns:
//   - A pointer to a Response[T] containing the provided status code, body, and error message.
func NewResponse[T any](statusCode int, body T, errorMsg ErrorType) *Response[T] {
	return &Response[T]{
		StatusCode: statusCode,
		Body:       body,
		Error:      errorMsg,
	}
}

// NewSuccessResponse creates a new Response with HTTP status 200 (OK) and the provided body.
// The error message is set to an empty string.
// T represents the type of the response body.
func NewSuccessResponse[T any](body T) *Response[T] {
	return NewResponse(http.StatusOK, body, "")
}
