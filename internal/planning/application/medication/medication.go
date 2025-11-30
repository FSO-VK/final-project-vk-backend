package medication

import "github.com/google/uuid"

type MedicationService interface {
	MedicationName(id uuid.UUID) (string, error)
}
