package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a list of fields which failed validation.
type FieldError struct {
	Field string
	Tag   string
	Value any
}

// ValidationError is a custom validator error.
type ValidationError struct {
	// description is a description of validation error
	description string

	// fields is a map of field name which failed validation
	// and the failed tag of rule.
	fields map[string]*FieldError
}

// Error is a method to implement error interface.
func (e *ValidationError) Error() string {
	if len(e.fields) == 0 {
		return e.description
	}

	fieldWithErrors := make([]string, 0)
	for name := range e.fields {
		fieldWithErrors = append(fieldWithErrors, name)
	}

	return fmt.Sprintf("%s for fields: %s", e.description, fieldWithErrors)
}

var ErrCantValidate = &ValidationError{
	description: "type can't be validated",
}

// Fields returns a map of fields which failed validation.
func (e *ValidationError) Fields() map[string]*FieldError {
	return e.fields
}

// Validator is a validator interface.
type Validator interface {
	ValidateStruct(v any) *ValidationError
}

// ValidationProvider is an implementation of Validator.
type ValidationProvider struct {
	validate *validator.Validate
}

// NewValidationProvider returns a new instance of ValidationProvider.
func NewValidationProvider() *ValidationProvider {
	return &ValidationProvider{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// ValidateStruct is a method to validate given struct.
func (va *ValidationProvider) ValidateStruct(v any) *ValidationError {
	return va.errorMapping(va.validate.Struct(v))
}

// errorMapping is a helper to map library errors to custom errors.
func (va *ValidationProvider) errorMapping(err error) *ValidationError {
	if err == nil {
		return nil
	}

	// Error can't be wrapped cause it's a library error
	//nolint:errorlint
	_, ok := err.(*validator.InvalidValidationError)
	if ok {
		return ErrCantValidate
	}

	fields := make(map[string]*FieldError)
	// Error can't be wrapped cause it's a library error
	//nolint:errorlint
	fieldErrors, _ := err.(validator.ValidationErrors)
	for _, fieldErr := range fieldErrors {
		fields[fieldErr.Field()] = &FieldError{
			Field: fieldErr.Field(),
			Tag:   fieldErr.Tag(),
			Value: fieldErr.Value(),
		}
	}

	return &ValidationError{
		description: "validation failed",
		fields:      fields,
	}
}
