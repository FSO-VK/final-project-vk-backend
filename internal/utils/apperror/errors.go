package apperror

import "fmt"

type category string

const (
	UserCategory     category = "user error"
	SystemCategory   category = "system error"
	InternalCategory category = "internal error"
)

type ApplicationError struct {
	cat  category
	name string
	err  error
}

// Error implements error interface
func (ae *ApplicationError) Error() string {
	return fmt.Sprintf("%s: %s: %s", ae.cat, ae.name, ae.err.Error())
}

func (ae *ApplicationError) Unwrap() error {
	return ae.err
}

func (ae *ApplicationError) Name() string {
	return ae.name
}

func (ae *ApplicationError) Category() category {
	return ae.cat
}

func User(name string, err error) *ApplicationError {
	return &ApplicationError{
		cat:  UserCategory,
		name: name,
		err: err,
	}
}

func System(name string, err error) *ApplicationError {
	return &ApplicationError{
		cat:  SystemCategory,
		name: name,
		err: err,
	}
}

func Internal(name string, err error) *ApplicationError {
	return &ApplicationError{
		cat:  InternalCategory,
		name: name,
		err: err,
	}
}