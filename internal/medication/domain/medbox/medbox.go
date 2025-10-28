// Package medbox implements domain layer for medication box aggregate.
package medbox

import (
	"github.com/google/uuid"
)

// MedicationBox is a domain entity that represents a user's medication box aggregate root.
type MedicationBox struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	MedicationsID []uuid.UUID
}
