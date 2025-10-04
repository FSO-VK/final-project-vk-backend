package medicine

import (
	"context"
)

type MedicineRepository interface {
	Create(ctx context.Context, medicine *Medicine) (*Medicine, error)
	GetByID(ctx context.Context, medicineID uint) (*Medicine, error)
	GetListByID(ctx context.Context, medicinesID []uint) ([]*Medicine, error)
	GetListAll(ctx context.Context) ([]*Medicine, error)
	Update(ctx context.Context, medicine *Medicine) (*Medicine, error)
	Delete(ctx context.Context, medicineID uint) error
}
