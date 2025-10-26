package memory

import (
	"context"
	"strconv"
	"sync"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
)

// MedicationStorage is a storage for medications.
type MedicationStorage struct {
	data  *Cache[*medication.Medication]
	count uint

	mu *sync.RWMutex
}

// NewMedicationStorage returns a new MedicationStorage.
func NewMedicationStorage() *MedicationStorage {
	return &MedicationStorage{
		data:  NewCache[*medication.Medication](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Create creates a new medication in memory.
func (s *MedicationStorage) Create(
	_ context.Context,
	medication *medication.Medication,
) (*medication.Medication, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.count++
	s.data.Set(medication.ID.String(), medication)
	return medication, nil
}

// GetByID returns a medication by id.
func (s *MedicationStorage) GetByID(
	_ context.Context,
	medicationID uint,
) (*medication.Medication, error) {
	drug, ok := s.data.Get(strconv.FormatUint(uint64(medicationID), 10))
	if !ok {
		return nil, medication.ErrNoMedicationFound
	}
	return drug, nil
}

// GetListAll returns a list of all medications.
func (s *MedicationStorage) GetListAll(_ context.Context) ([]*medication.Medication, error) {
	list := make([]*medication.Medication, 0)
	for _, medication := range s.data.data {
		list = append(list, medication)
	}
	return list, nil
}

// Update updates a medication in memory.
func (s *MedicationStorage) Update(
	_ context.Context,
	medicationToUpdate *medication.Medication,
) (*medication.Medication, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(medicationToUpdate.ID.String())
	if !ok {
		s.data.mu.Unlock()
		return nil, medication.ErrNoMedicationFound
	}
	s.data.Set(medicationToUpdate.ID.String(), medicationToUpdate)
	return medicationToUpdate, nil
}

// Delete deletes a medication in memory.
func (s *MedicationStorage) Delete(_ context.Context, medicationID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(strconv.FormatUint(uint64(medicationID), 10))
	if !ok {
		return medication.ErrNoMedicationFound
	}

	s.data.Delete(strconv.FormatUint(uint64(medicationID), 10))
	return nil
}
