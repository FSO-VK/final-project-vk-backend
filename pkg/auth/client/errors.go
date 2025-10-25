package client

import "errors"

var (
	ErrUnauthorized           = errors.New("http-auth: unauthorized")
	ErrAuthServiceUnavailable = errors.New("http-auth: service unavailable")
	ErrInvalidAuthResponse    = errors.New("http-auth: invalid response")
	ErrMarshalingRequest      = errors.New("http-auth: marshaling error")
	ErrBadResponse            = errors.New("http-auth: unexpected status code in response")
)
