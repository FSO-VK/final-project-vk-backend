package medication

import (
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const maxActiveSubstanceLength = 200

var ErrInvalidActiveSubstance = errors.New("invalid active substance")

// ActiveSubstanceDraft represents a active substance draft entity.
type ActiveSubstanceDraft struct {
	Name  string
	Value float32
	Unit  string
}

// ActiveSubstance is a VO representing the active
// pharmaceutical substance and its dosage.
type ActiveSubstance struct {
	name string
	dose Amount
}

// NewActiveSubstance creates validated medication active substance.
func NewActiveSubstance(
	activeSubstance ActiveSubstanceDraft,
) (ActiveSubstance, error) {
	dose, err := NewAmount(AmountDraft{
		Value: activeSubstance.Value,
		Unit:  activeSubstance.Unit})

	err = errors.Join(
		validation.MaxLength(activeSubstance.Name, maxActiveSubstanceLength),
		err,
	)
	if err != nil {
		return ActiveSubstance{}, fmt.Errorf("%w: %w", ErrInvalidActiveSubstance, err)
	}

	return ActiveSubstance{
		name: activeSubstance.Name,
		dose: dose,
	}, nil
}

// Name returns the name of the active substance.
func (a ActiveSubstance) Name() string { return a.name }

// Dose returns the dose amount of the active substance.
func (a ActiveSubstance) Dose() Amount { return a.dose }
