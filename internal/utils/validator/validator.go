package validator

type Validator interface {
	ValidateStruct(v any) error
}
