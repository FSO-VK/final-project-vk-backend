package medicine

import (
	"context"
	"errors"
)

// ErrNoMedicineFound is an error when a medicine is not found.
var ErrNoMedicineFound = errors.New("medicine not found")

// RepositoryForMedication is a repository - provides methods to create, get, update, and delete medicines.
type RepositoryForMedication interface {
	Create(ctx context.Context, medicine *Medicine) (*Medicine, error)
	GetByID(ctx context.Context, medicineID uint) (*Medicine, error)
	GetListAll(ctx context.Context) ([]*Medicine, error)
	Update(ctx context.Context, medicine *Medicine) (*Medicine, error)
	Delete(ctx context.Context, medicineID uint) error
}
