package client

import "errors"

var (
	// ErrBadResponse is returned when the response from authorized is not 200 OK.
	ErrAuthServiceUnavailable = errors.New("http-auth: service unavailable")
	// ErrInvalidAuthResponse is returned when the response is not like expected.
	ErrInvalidAuthResponse = errors.New("http-auth: invalid response")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("http-auth: invalid request")
	// ErrMarshalingRequest is returned when the request cannot be marshaled.
	ErrMarshalingRequest = errors.New("http-auth: marshaling error")
	// ErrBadResponse is returned when the status code in response is not 200 OK.
	ErrBadResponse = errors.New("http-auth: unexpected status code in response")
)
