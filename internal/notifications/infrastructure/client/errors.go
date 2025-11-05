package client

import "errors"

var (
	// ErrPushServiceUnavailable is returned when the response from push endpoint is not 200 OK.
	ErrPushServiceUnavailable = errors.New("notifications push api: service unavailable")
	// ErrBadRequest is returned when the request is invalid.
	ErrBadRequest = errors.New("notifications push api: bad request")
)
