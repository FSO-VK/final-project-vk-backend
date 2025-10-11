package validator

type Validator interface {
	ValidateStruct(v any) error
}

type Validation struct {
}

func NewValidation() *Validation {
	return &Validation{}
}

func (val *Validation) ValidateStruct(v any) error {
	return nil
}
