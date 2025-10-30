package memory

import (
	"context"
	"sync"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/google/uuid"
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
	s.data.Set(medication.GetID().String(), medication)
	return medication, nil
}

// GetByID returns a medication by id.
func (s *MedicationStorage) GetByID(
	_ context.Context,
	medicationID uuid.UUID,
) (*medication.Medication, error) {
	drug, ok := s.data.Get(medicationID.String())
	if !ok {
		return nil, medication.ErrNoMedicationFound
	}
	return drug, nil
}

// Update updates a medication in memory.
func (s *MedicationStorage) Update(
	_ context.Context,
	medicationToUpdate *medication.Medication,
) (*medication.Medication, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(medicationToUpdate.GetID().String())
	if !ok {
		s.data.mu.Unlock()
		return nil, medication.ErrNoMedicationFound
	}
	s.data.Set(medicationToUpdate.GetID().String(), medicationToUpdate)
	return medicationToUpdate, nil
}

// Delete deletes a medication in memory.
func (s *MedicationStorage) Delete(_ context.Context, medicationID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Get(medicationID.String())
	if !ok {
		return medication.ErrNoMedicationFound
	}

	s.data.Delete(medicationID.String())
	return nil
}
