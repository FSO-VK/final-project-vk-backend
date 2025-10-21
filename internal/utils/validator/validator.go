package validator

import "github.com/go-playground/validator/v10"

// FieldErrors represents a list of fields which failed validation.
type FieldErrors struct {
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
	fields map[string]*FieldErrors
}

// Error is a method to implement error interface.
func (e *ValidationError) Error() string {
	if len(e.fields) == 0 {
		return e.description
	}
	return "validation error"
}

var (
	ErrCantValidate = &ValidationError{
		description: "type can't be validated",
	}
)

// Fields returns a map of fields which failed validation.
func (e *ValidationError) Fields() map[string]*FieldErrors {
	return e.fields
}

// Validator is a validator interface.
type Validator interface {
	ValidateStruct(v any) error
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
func (va *ValidationProvider) ValidateStruct(v any) error {
	return va.errorMapping(va.validate.Struct(v))
}

// errorMapping is a helper to map library errors to custom errors.
func (va *ValidationProvider) errorMapping(err error) *ValidationError {
	if err == nil {
		return nil
	}

	_, ok := err.(*validator.InvalidValidationError)
	if ok {
		return ErrCantValidate
	}

	fields := make(map[string]*FieldErrors)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		fields[fieldErr.Field()] = &FieldErrors{
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