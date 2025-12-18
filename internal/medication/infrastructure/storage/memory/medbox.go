package memory

import (
	"context"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

// MedicationBoxStorage is a storage for MedicationBoxes.
type MedicationBoxStorage struct {
	data  *cache.Cache[*medbox.MedicationBox]
	count uint

	mu *sync.RWMutex
}

// NewMedicationBoxStorage returns a new MedicationBoxStorage.
func NewMedicationBoxStorage() *MedicationBoxStorage {
	return &MedicationBoxStorage{
		data:  cache.NewCache[*medbox.MedicationBox](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// SetMedicationBox creates a new MedicationBox in memory.
func (s *MedicationBoxStorage) SetMedicationBox(
	_ context.Context,
	medicationBox *medbox.MedicationBox,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if medicationBox == nil {
		return medbox.ErrNoMedicationBoxFound
	}
	_, ok := s.data.Get(medicationBox.GetID().String())
	if !ok {
		return medbox.ErrNoMedicationBoxFound
	}
	s.data.Set(medicationBox.GetID().String(), medicationBox)
	return nil
}

// GetMedicationBox returns a medication box by userID.
func (s *MedicationBoxStorage) GetMedicationBox(
	_ context.Context,
	userID uuid.UUID,
) (*medbox.MedicationBox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, medicationBox := range s.data.GetAll() {
		if medicationBox.GetUserID() == userID {
			return medicationBox, nil
		}
	}
	return nil, medbox.ErrNoMedicationBoxFound
}

// CreateMedicationBox creates a new medication in memory.
func (s *MedicationBoxStorage) CreateMedicationBox(
	_ context.Context,
	medicationBox *medbox.MedicationBox,
) (*medbox.MedicationBox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.count++
	s.data.Set(medicationBox.GetID().String(), medicationBox)
	return medicationBox, nil
}
