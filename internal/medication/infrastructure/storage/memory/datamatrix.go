package memory

import (
	"context"
	"sync"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
)

// DataMatrixStorage is a storage for medications.
type DataMatrixStorage struct {
	data  *Cache[*medication.Medication]
	count uint
	mu    *sync.RWMutex
}

// NewDataMatrixStorage returns a new DataMatrixStorage.
func NewDataMatrixStorage() *DataMatrixStorage {
	return &DataMatrixStorage{
		data: NewCache[*medication.Medication](),
		count:   0,
		mu:   &sync.RWMutex{},
	}
}

// Create creates a new medication in memory.
func (s *DataMatrixStorage) Create(
	_ context.Context,
	dataMatrixString string,
	medication *medication.Medication,
) (*medication.Medication, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.count++
	s.data.Set(dataMatrixString, medication)
	return medication, nil
}

// GetByDataMatrixString returns a medication info by GTIN+SerialNumber+CryptoData91+CryptoData92.
func (s *DataMatrixStorage) GetByDataMatrixString(
	_ context.Context,
	dataMatrixString string,
) (*medication.Medication, error) {
	drug, ok := s.data.Get(dataMatrixString)
	if !ok {
		return nil, medication.ErrNoMedicationFound
	}
	return drug, nil
}
