package medication

import (
	"errors"
	"fmt"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const (
	maxNameLength              = 200
	maxInternationalNameLength = 100
)

var (
	ErrInvalidName              = errors.New("invalid name")
	ErrInvalidInternationalName = errors.New("invalid international name")
)

// Name is a Value Object representing the name of a medication.
type Name string

// NewName creates validated medication name.
func NewName(name string) (Name, error) {
	err := errors.Join(
		validation.Required(name),
		validation.MaxLength(name, maxNameLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidName, err)
	}
	return Name(name), nil
}

// String implements stringer interface for Name.
func (n Name) String() string {
	return string(n)
}

// InternationalName is a Value Object representing the international non-proprietary name of a medication.
type InternationalName string

// NewInternationalName creates validated medication international name.
func NewInternationalName(name string) (InternationalName, error) {
	err := errors.Join(
		validation.MaxLength(name, maxInternationalNameLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidInternationalName, err)
	}
	return InternationalName(name), nil
}

// String implements stringer interface for InternationalName.
func (n InternationalName) String() string {
	return string(n)
}
