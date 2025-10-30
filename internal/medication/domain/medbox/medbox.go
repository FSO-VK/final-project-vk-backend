// Package medbox implements domain layer for medication box aggregate.
package medbox

import (
	"slices"

	"github.com/google/uuid"
)

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

// RemoveMedication removes a medication from the medication box.
func (m *MedicationBox) RemoveMedication(medicationID uuid.UUID) {
	m.medicationsID = slices.DeleteFunc(m.medicationsID, func(id uuid.UUID) bool {
		return id == medicationID
	})
}
