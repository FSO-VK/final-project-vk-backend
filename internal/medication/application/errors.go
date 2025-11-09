package application

import "errors"

// Common errors for application layer.
var (
	// ErrValidationFail indicates that struct validation failed.
	ErrValidationFail = errors.New("struct validation failed")
	// ErrNoMedication indicates that no medication was found.
	ErrNoMedication = errors.New("no medication")
)
