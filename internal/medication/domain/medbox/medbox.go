// Package medbox implements domain layer for medication box aggregate.
package medbox

import (
	"errors"
	"slices"

	"github.com/google/uuid"
)

// ErrNoMedication indicates that medication is not found in medication box.
var ErrNoMedication = errors.New("no medication in medication box")

// MedicationBox is a domain entity that represents a user's medication box aggregate root.
type MedicationBox struct {
	id            uuid.UUID
	userID        uuid.UUID
	medicationsID []uuid.UUID
}

// NewMedicationBox creates a new medication box aggregate root.
func NewMedicationBox(userID uuid.UUID) *MedicationBox {
	return &MedicationBox{
		id:            uuid.New(),
		userID:        userID,
		medicationsID: []uuid.UUID{},
	}
}

// GetID returns the unique identifier of the medication box.
func (m *MedicationBox) GetID() uuid.UUID {
	return m.id
}

// GetUserID returns the unique identifier of the user.
func (m *MedicationBox) GetUserID() uuid.UUID {
	return m.userID
}

// GetMedicationsID returns the unique identifiers of the medications.
func (m *MedicationBox) GetMedicationsID() []uuid.UUID {
	return m.medicationsID
}

// AddMedication adds a medication to the medication box.
func (m *MedicationBox) AddMedication(medicationID uuid.UUID) {
	m.medicationsID = append(m.medicationsID, medicationID)
}

// HasMedication checks if a medication is in the medication box.
func (m *MedicationBox) HasMedication(medicationID uuid.UUID) bool {
	return slices.ContainsFunc(m.medicationsID, func(id uuid.UUID) bool {
		return id == medicationID
	})
}

// RemoveMedication removes a medication from the medication box.
func (m *MedicationBox) RemoveMedication(medicationID uuid.UUID) error {
	if !m.HasMedication(medicationID) {
		return ErrNoMedication
	}
	m.medicationsID = slices.DeleteFunc(m.medicationsID, func(id uuid.UUID) bool {
		return id == medicationID
	})
	return nil
}
