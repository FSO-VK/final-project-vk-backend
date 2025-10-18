package medication

import (
	"context"
	"errors"
)

// ErrNoMedicationFound is an error when a medication is not found.
var ErrNoMedicationFound = errors.New("medication not found")

// RepositoryForMedication is a domain repository interface that defines
// data access contract for medication aggregate.
type RepositoryForMedication interface {
	Create(ctx context.Context, medication *Medication) (*Medication, error)
	GetByID(ctx context.Context, medicationID uint) (*Medication, error)
	GetListAll(ctx context.Context) ([]*Medication, error)
	Update(ctx context.Context, medication *Medication) (*Medication, error)
	Delete(ctx context.Context, medicationID uint) error
}
