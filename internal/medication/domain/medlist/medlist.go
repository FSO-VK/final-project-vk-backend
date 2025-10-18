// Package medlist implements domain layer for medication list aggregate.
package medlist

import (
	"github.com/google/uuid"
)

// MedicationList is a domain entity that represents a user's medication list aggregate root.
type MedicationList struct {
	ID            uint
	UserID        uuid.UUID
	MedicationsID []uint
}
