package memory

import (
	"context"
	"strconv"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
)

// MedicineStorage is a storage for medicines.
type MedicineStorage struct {
	data *Cache[*medicine.Medicine]
	id   uint

	mu *sync.RWMutex
}

// NewMedicineStorage returns a new MedicineStorage.
func NewMedicineStorage() *MedicineStorage {
	return &MedicineStorage{
		data: NewCache[*medicine.Medicine](),
		id:   0,
		mu:   &sync.RWMutex{},
	}
}

// Create creates a new medicine in memory.
func (s *MedicineStorage) Create(
	_ context.Context,
	medicine *medicine.Medicine,
) (*medicine.Medicine, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	medicine.ID = s.id
	s.id++
	s.data.Set(strconv.FormatUint(uint64(medicine.ID), 10), medicine)
	return medicine, nil
}

// GetByID returns a medicine by id.
func (s *MedicineStorage) GetByID(
	_ context.Context,
	medicineID uint,
) (*medicine.Medicine, error) {
	drug, ok := s.data.Get(strconv.FormatUint(uint64(medicineID), 10))
	if !ok {
		return nil, medicine.ErrNoMedicineFound
	}
	return drug, nil
}

// GetListAll returns a list of all medicines.
func (s *MedicineStorage) GetListAll(_ context.Context) ([]*medicine.Medicine, error) {
	list := make([]*medicine.Medicine, 0)
	for _, medicine := range s.data.data {
		list = append(list, medicine)
	}
	return list, nil
}

// Update updates a medicine in memory.
func (s *MedicineStorage) Update(
	_ context.Context,
	medicineToUpdate *medicine.Medicine,
) (*medicine.Medicine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(strconv.FormatUint(uint64(medicineToUpdate.ID), 10))
	if !ok {
		s.data.mu.Unlock()
		return nil, medicine.ErrNoMedicineFound
	}
	s.data.Set(strconv.FormatUint(uint64(medicineToUpdate.ID), 10), medicineToUpdate)
	return medicineToUpdate, nil
}

// Delete deletes a medicine in memory.
func (s *MedicineStorage) Delete(_ context.Context, medicineID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(strconv.FormatUint(uint64(medicineID), 10))
	if !ok {
		return medicine.ErrNoMedicineFound
	}

	s.data.Delete(strconv.FormatUint(uint64(medicineID), 10))
	return nil
}
