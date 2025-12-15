package medication

import "github.com/google/uuid"

type MedicationService interface {
	MedicationName(id uuid.UUID, userID uuid.UUID) (string, error)
}
