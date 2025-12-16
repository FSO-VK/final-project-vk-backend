package medication

import "github.com/google/uuid"

// MedicationService is an interface for getting medication info from medication service.
type MedicationService interface {
	MedicationName(id uuid.UUID, userID uuid.UUID) (string, error)
}
