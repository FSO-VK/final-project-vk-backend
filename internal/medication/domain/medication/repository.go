package medication

import (
	"context"
	"errors"
	"iter"
	"time"

	"github.com/google/uuid"
)

// ErrNoMedicationFound is an error when a medication is not found.
var ErrNoMedicationFound = errors.New("medication not found")

// Repository is a domain repository interface that defines
// data access contract for medication aggregate.
type Repository interface {
	Create(ctx context.Context, medication *Medication) (*Medication, error)
	GetByID(ctx context.Context, medicationID uuid.UUID) (*Medication, error)
	Update(ctx context.Context, medication *Medication) (*Medication, error)
	Delete(ctx context.Context, medicationID uuid.UUID) error
	MedicationByExpiration(
		ctx context.Context,
		timeDelta time.Duration,
	) (iter.Seq[*Medication], error)
}
