package client

import "errors"

var (
	ErrAuthServiceUnavailable = errors.New("http-auth: service unavailable")
	ErrInvalidAuthResponse    = errors.New("http-auth: invalid response")
	ErrInvalidREquest         = errors.New("http-auth: invalid request")
	ErrMarshalingRequest      = errors.New("http-auth: marshaling error")
	ErrBadResponse            = errors.New("http-auth: unexpected status code in response")
)
