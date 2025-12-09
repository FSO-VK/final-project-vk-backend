package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

// errGotNilIntakeRecord is an error when save gets nil intake record to add.
var errGotNilIntakeRecord = errors.New("cannot save nil intake record")

// RecordStorage is a storage for Records.
type RecordStorage struct {
	data  *cache.Cache[*record.IntakeRecord]
	count uint

	mu *sync.RWMutex
}

// NewRecordStorage returns a new RecordStorage.
func NewRecordStorage() *RecordStorage {
	return &RecordStorage{
		data:  cache.NewCache[*record.IntakeRecord](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Create creates a new record in memory.
func (s *RecordStorage) Save(
	_ context.Context,
	newRecord *record.IntakeRecord,
) error {
	if newRecord == nil {
		return errGotNilIntakeRecord
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	s.data.Set(newRecord.ID().String(), newRecord)
	return nil
}

// Create creates a bulk of new records in memory.
func (s *RecordStorage) SaveBulk(
	_ context.Context,
	bulkOfRecords []*record.IntakeRecord,
) error {
	if bulkOfRecords == nil {
		return errGotNilIntakeRecord
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, oneRecord := range bulkOfRecords {
		if oneRecord == nil {
			return errGotNilIntakeRecord
		}
		s.count++
		s.data.Set(oneRecord.ID().String(), oneRecord)
	}
	return nil
}

// GetByID returns a record by id.
func (s *RecordStorage) GetByID(
	_ context.Context,
	id uuid.UUID,
) (*record.IntakeRecord, error) {
	requestedRecord, ok := s.data.Get(id.String())
	if !ok {
		return nil, record.ErrNoRecordFound
	}
	return requestedRecord, nil
}

// UserRecords returns all records by plan id.
func (s *RecordStorage) GetByPlanID(
	_ context.Context,
	userID uuid.UUID,
) ([]*record.IntakeRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*record.IntakeRecord
	for _, oneRecord := range s.data.GetAll() {
		if oneRecord.PlanID() == userID {
			result = append(result, oneRecord)
		}
	}

	return result, nil
}
