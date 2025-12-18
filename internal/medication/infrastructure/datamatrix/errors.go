package datamatrix

import "errors"

var (
	// ErrAuthServiceUnavailable is returned when the response from dataMatrix api is not 200 OK.
	ErrAuthServiceUnavailable = errors.New("dataMatrix api: service unavailable")
	// ErrNoMedicationFound is returned when the response is not like expected ("codeFounded":false).
	ErrNoMedicationFound = errors.New("dataMatrix api: no medication found")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("dataMatrix api: invalid request")
	// ErrBadResponse is returned when the response status code is not 200 or body is not like expected.
	ErrBadResponse = errors.New(
		"dataMatrix api: invalid response not 200 or body is not like expected",
	)
)
