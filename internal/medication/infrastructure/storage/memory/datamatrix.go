package memory

import (
	"context"
	"strconv"
	"sync"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
)

// DataMatrixStorage is a storage for medications.
type DataMatrixStorage struct {
	data *Cache[*medication.Medication]
	id   uint
	mu *sync.RWMutex
}

// NewDataMatrixStorage returns a new DataMatrixStorage.
func NewDataMatrixStorage() *DataMatrixStorage {
	return &DataMatrixStorage{
		data: NewCache[*medication.Medication](),
		id:   0,
		mu:   &sync.RWMutex{},
	}
}

// Create creates a new medication info from API in cache.
func (s *DataMatrixStorage) Create(
	_ context.Context,
	dataMatrixString string,
	medication *medication.Medication,
) (*medication.Medication, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	medication.ID = s.id
	s.id++
	s.data.Set(strconv.FormatUint(uint64(medication.ID), 10), medication)
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
