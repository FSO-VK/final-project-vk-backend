package notificationclient

import "errors"

var (
	// ErrNotificationServiceUnavailable is returned when the response from Notification api is unavailable.
	ErrNotificationServiceUnavailable = errors.New("notification api: service unavailable")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("notification api: invalid request")
	// ErrBadResponse is returned when the response status code is not 200 or body is not like expected.
	ErrBadResponse = errors.New(
		"notification api: invalid response not 200 or body is not like expected",
	)
)
