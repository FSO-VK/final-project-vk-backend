package api

import "net/http"

type Response[T any] struct {
	StatusCode int       `json:"statusCode"`
	Body       T         `json:"body"`
	Error      ErrorType `json:"error,omitempty"`
}

func NewResponse[T any](statusCode int, body T, errorMsg ErrorType) *Response[T] {
	return &Response[T]{
		StatusCode: statusCode,
		Body:       body,
		Error:      errorMsg,
	}
}

func NewSuccessResponse[T any](body T) *Response[T] {
	return NewResponse(http.StatusOK, body, "")
}
