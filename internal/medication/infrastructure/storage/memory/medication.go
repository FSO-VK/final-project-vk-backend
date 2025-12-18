package memory

import (
	"context"
	"iter"
	"sync"
	"time"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

// MedicationStorage is a storage for medications.
type MedicationStorage struct {
	data  *cache.Cache[*medication.Medication]
	count uint

	mu *sync.RWMutex
}

// NewMedicationStorage returns a new MedicationStorage.
func NewMedicationStorage() *MedicationStorage {
	return &MedicationStorage{
		data:  cache.NewCache[*medication.Medication](),
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

// MedicationByExpiration returns all medications expiring within timeDelta from now.
func (s *MedicationStorage) MedicationByExpiration(
	_ context.Context,
	timeDelta time.Duration,
) (iter.Seq[*medication.Medication], error) {
	return func(yield func(*medication.Medication) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		all := s.data.GetAll()
		now := time.Now()
		expirationThreshold := now.Add(timeDelta)

		for _, med := range all {
			exp := med.GetExpirationDate()
			if !exp.IsZero() &&
				!exp.Before(now) &&
				!exp.After(expirationThreshold) {
				if !yield(med) {
					return
				}
			}
		}
	}, nil
}
