package medicine

import (
	"context"
	"errors"
)

var ErrNoMedicineFound = errors.New("medicine not found")

type MedicineRepository interface {
	Create(ctx context.Context, medicine *Medicine) (*Medicine, error)
	GetByID(ctx context.Context, medicineID uint) (*Medicine, error)
	GetListAll(ctx context.Context) ([]*Medicine, error)
	Update(ctx context.Context, medicine *Medicine) (*Medicine, error)
	Delete(ctx context.Context, medicineID uint) error
}
