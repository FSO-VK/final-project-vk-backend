package medbox

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNoMedicationBoxFound is an error when a medication box is not found.
var ErrNoMedicationBoxFound = errors.New("medication box not found")

// RepositoryForMedicationBox is a domain repository interface that defines
// data access contract for medication box aggregate.
type RepositoryForMedicationBox interface {
	SetMedicationBox(ctx context.Context, medicationBox *MedicationBox) error
	CreateMedicationBox(ctx context.Context, UserID uuid.UUID) (*MedicationBox, error)
	GetMedicationBox(ctx context.Context, UserID uuid.UUID) (*MedicationBox, error)
}
