package medication

import (
	"errors"
	"fmt"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

var (
	ErrInvalidAmount = errors.New("invalid amount")
)

// AmountDraft is a helper struct for constructing
// Amount.
type AmountDraft struct {
	Value float32
	Unit  string
}

// Amount is a VO representing the quantity of a medication
// with its unit of measurement.
type Amount struct {
	value float32
	unit  Unit
}

// NewAmount creates validated medication amount.
func NewAmount(draft AmountDraft) (Amount, error) {
	medicationUnit, err := NewUnit(draft.Unit)

	err = errors.Join(
		err,
		validation.Positive(draft.Value),
	)
	if err != nil {
		return Amount{}, fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	return Amount{
		value: draft.Value,
		unit:  medicationUnit,
	}, nil
}

// Value returns the numeric value of the amount.
func (a Amount) Value() float32 {
	return a.value
}

// Unit returns the unit of measurement for the amount.
func (a Amount) Unit() Unit {
	return a.unit
}
