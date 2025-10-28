package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/google/uuid"
)

// MedicationBoxStorage is a storage for MedicationBoxes.
type MedicationBoxStorage struct {
	data  *Cache[*medbox.MedicationBox]
	count uint

	mu *sync.RWMutex
}

// NewMedicationBoxStorage returns a new MedicationBoxStorage.
func NewMedicationBoxStorage() *MedicationBoxStorage {
	return &MedicationBoxStorage{
		data:  NewCache[*medbox.MedicationBox](),
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
	_, ok := s.data.Get(medicationBox.ID.String())
	if !ok {
		s.data.mu.Unlock()
		return medbox.ErrNoMedicationBoxFound
	}
	fmt.Println("11111111111111111111 written")
	s.data.Set(medicationBox.ID.String(), medicationBox)
	return nil
}

// GetMedicationBox returns a medication box by userID.
func (s *MedicationBoxStorage) GetMedicationBox(
	_ context.Context,
	userID uuid.UUID,
) (*medbox.MedicationBox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, medicationBox := range s.data.data {
		_ = key
		if medicationBox.UserID == userID {
			return medicationBox, nil
		}
	}
	return nil, medbox.ErrNoMedicationBoxFound
}

// Create creates a new medication in memory.
func (s *MedicationBoxStorage) CreateMedicationBox(
	_ context.Context,
	userID uuid.UUID,
) (*medbox.MedicationBox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	medicationBoxID := uuid.New()

	medicationBox := &medbox.MedicationBox{
		ID:            medicationBoxID,
		UserID:        userID,
		MedicationsID: []uuid.UUID{},
	}
	s.count++
	s.data.Set(medicationBoxID.String(), medicationBox)
	return medicationBox, nil
}
