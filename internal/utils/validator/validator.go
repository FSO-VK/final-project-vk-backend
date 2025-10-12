package validator

type Validator interface {
	ValidateStruct(v any) error
}

type ValidationProvider struct{}

func NewValidationProvider() *ValidationProvider {
	return &ValidationProvider{}
}

func (val *ValidationProvider) ValidateStruct(v any) error {
	return nil
}
