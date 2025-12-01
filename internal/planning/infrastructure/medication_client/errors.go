package medicationclient

import "errors"

var (
	// ErrMedicationServiceUnavailable is returned when the response from medication api is unavailable.
	ErrMedicationServiceUnavailable = errors.New("medication api: service unavailable")
	// ErrNoMedicationFound is returned when the response is not like expected ("codeFounded":false).
	ErrNoMedicationFound = errors.New("medication api: no medication found")
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("medication api: invalid request")
	// ErrBadResponse is returned when the response status code is not 200 or body is not like expected.
	ErrBadResponse = errors.New(
		"medication api: invalid response not 200 or body is not like expected",
	)
)
