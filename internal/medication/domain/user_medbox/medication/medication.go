// Package medication is a domain for medication
package medication

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// MedicationDraft represents a medication draft entity
// that uses built-in types.
type MedicationDraft struct {
	ID         uuid.UUID
	Info       MedicationInfoDraft
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DataMatrix string
}

// Medication represents a medication entity.
type Medication struct {
	id         uuid.UUID
	info       MedicationInfo
	createdAt  time.Time
	updatedAt  time.Time
	dataMatrix DataMatrix
}

// NewMedication creates a new medication.
func NewMedication(
	d MedicationDraft,
) (Medication, error) {
	var allErrors error
	info, err := NewMedicationInfo(d.Info)
	allErrors = errors.Join(allErrors, err)

	var dataMatrix DataMatrix
	if d.DataMatrix == "" {
		dataMatrix = NewEmptyDataMatrix()
	} else {
		dataMatrix, err = ParseDataMatrix(d.DataMatrix)
		allErrors = errors.Join(allErrors, err)
	}

	if allErrors != nil {
		return Medication{}, err
	}

	return Medication{
		id:         d.ID,
		info:       info,
		createdAt:  d.CreatedAt,
		updatedAt:  d.UpdatedAt,
		dataMatrix: dataMatrix,
	}, nil

}

// UpdateInfo updates medication info.
func (m *Medication) UpdateInfo(info MedicationInfo) {
	m.info = info
	m.updatedAt = time.Now()
}

func (m *Medication) Take(amount float32) error {
	// Since info is a copy of a real info, changing amount is safe
	// It is a bad code, but safe enough to ignore
	info := m.Info()
	currentAmountValue := info.Amount().value
	newAmountValue := max(currentAmountValue-amount, 0)

	newAmount, err := NewAmount(
		AmountDraft{
			Value: newAmountValue,
			Unit:  info.Amount().Unit().String(),
		},
	)
	if err != nil {
		// basically impossible but the linter must be shutted up
		return err
	}

	info.amount = newAmount

	m.UpdateInfo(info)
	return nil
}

// IsRanOut checks if amount of medication is zero
func (m *Medication) IsRanOut() bool {
	return m.info.amount.value == 0
}

// ID returns the unique identifier of the medication.
func (m *Medication) ID() uuid.UUID { return m.id }

// Info returns medication info.
func (m *Medication) Info() MedicationInfo {
	return m.info
}

// BarCode returns the DataMatrix of the medication.
func (m *Medication) DataMatrix() DataMatrix {
	return m.dataMatrix
}
