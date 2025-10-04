package memory

import (
	"context"
	"strconv"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medicine/medicine"
)

type MedicineStorage struct {
	data *Cache[*medicine.Medicine]
	id   uint

	mu *sync.RWMutex
}

func NewMedicineStorage() *MedicineStorage {
	return &MedicineStorage{
		data: NewCache[*medicine.Medicine](),
		id:   0,
		mu:   &sync.RWMutex{},
	}
}

func (s *MedicineStorage) Create(
	ctx context.Context,
	medicine *medicine.Medicine,
) (*medicine.Medicine, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	medicine.ID = s.id
	s.id++
	s.data.Set(strconv.Itoa(int(medicine.ID)), medicine)
	return medicine, nil
}

func (s *MedicineStorage) GetByID(
	ctx context.Context,
	medicineID uint,
) (*medicine.Medicine, error) {
	drug, ok := s.data.Get(strconv.Itoa(int(medicineID)))
	if !ok {
		return nil, medicine.ErrNoMedicineFound
	}
	return drug, nil
}

func (s *MedicineStorage) GetListAll(ctx context.Context) ([]*medicine.Medicine, error) {
	list := make([]*medicine.Medicine, 0)
	for _, medicine := range s.data.data {
		list = append(list, medicine)
	}
	return list, nil
}

func (s *MedicineStorage) Update(
	ctx context.Context,
	Medicine *medicine.Medicine,
) (*medicine.Medicine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(strconv.Itoa(int(Medicine.ID)))
	if !ok {
		s.data.mu.Unlock()
		return nil, medicine.ErrNoMedicineFound
	}
	s.data.Set(strconv.Itoa(int(Medicine.ID)), Medicine)
	return Medicine, nil
}

func (s *MedicineStorage) Delete(ctx context.Context, medicineID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(strconv.Itoa(int(medicineID)))
	if !ok {
		return medicine.ErrNoMedicineFound
	}

	s.data.Delete(strconv.Itoa(int(medicineID)))
	return nil
}
