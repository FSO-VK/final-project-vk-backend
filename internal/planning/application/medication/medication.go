package medication

import "github.com/google/uuid"

// MedicationService provides access to medication data.
type MedicationService interface {
	MedicationName(id uuid.UUID, userID uuid.UUID) (string, error)
}
