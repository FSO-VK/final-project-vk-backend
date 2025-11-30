package application

import "errors"

// Common errors for application layer.
var (
	// ErrValidationFail indicates that struct validation failed.
	ErrValidationFail = errors.New("struct validation failed")
	// ErrNoPlan indicates that no plan was found.
	ErrNoPlan = errors.New("no plan")
	// ErrNoIntakeRecord indicates that no intake record was found.
	ErrNoIntakeRecord = errors.New("no intake record")
	// ErrPlanNotFound is an error when plan is not belongs to user.
	ErrPlanNotBelongToUser = errors.New("plan does not belong to user")
)
