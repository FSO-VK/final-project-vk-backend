// Package memory contains cache for in memory db.
package memory

import (
	"context"
	"sync"

	client "github.com/FSO-VK/final-project-vk-backend/internal/medication/application/api_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
)

// DataMatrixStorage is a storage for medications.
type DataMatrixStorage struct {
	data  *cache.Cache[*client.MedicationInfo]
	count uint
	mu    *sync.RWMutex
}

// NewDataMatrixStorage returns a new DataMatrixStorage.
func NewDataMatrixStorage() *DataMatrixStorage {
	return &DataMatrixStorage{
		data:  cache.NewCache[*client.MedicationInfo](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Set creates a new medication in memory.
func (s *DataMatrixStorage) Set(
	_ context.Context,
	dataMatrixString string,
	medication *client.MedicationInfo,
) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.count++
	s.data.Set(dataMatrixString, medication)
	return nil
}

// Get returns a medication info by GTIN+SerialNumber+CryptoData91+CryptoData92.
func (s *DataMatrixStorage) Get(
	_ context.Context,
	dataMatrixString string,
) (*client.MedicationInfo, error) {
	drug, ok := s.data.Get(dataMatrixString)
	if !ok {
		return nil, client.ErrNoMedicationFound
	}
	return drug, nil
}
