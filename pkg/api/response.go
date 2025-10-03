// Package api provides common interaction protocol between frontend and backend
package api

import "net/http"

// Response is a generic struct representing a response from the API.
type Response[T any] struct {
	StatusCode int       `json:"statusCode"`
	Body       T         `json:"body"`
	Error      ErrorType `json:"error,omitempty"`
}

// NewResponse creates a new Response with the provided status code, body, and error message.
func NewResponse[T any](statusCode int, body T, errorMsg ErrorType) *Response[T] {
	return &Response[T]{
		StatusCode: statusCode,
		Body:       body,
		Error:      errorMsg,
	}
}

// NewSuccessResponse creates a new success response (with HTTP 200 OK) with the provided body.
func NewSuccessResponse[T any](body T) *Response[T] {
	return NewResponse(http.StatusOK, body, "")
}
