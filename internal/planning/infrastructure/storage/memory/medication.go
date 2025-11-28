package memory

import (
	"context"
	"sync"

	plan "github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

// PlanStorage is a storage for Plans.
type PlanStorage struct {
	data  *cache.Cache[*plan.Plan]
	count uint

	mu *sync.RWMutex
}

// NewPlanStorage returns a new PlanStorage.
func NewPlanStorage() *PlanStorage {
	return &PlanStorage{
		data:  cache.NewCache[*plan.Plan](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Create creates a new plan in memory.
func (s *PlanStorage) Save(_ context.Context, newPlan *plan.Plan) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.count++
	s.data.Set(newPlan.ID().String(), newPlan)
	return nil
}

// GetByID returns a plan by id.
func (s *PlanStorage) GetByID(
	_ context.Context,
	id uuid.UUID,
) (*plan.Plan, error) {
	requestedPlan, ok := s.data.Get(id.String())
	if !ok {
		return nil, plan.ErrNoPlanFound
	}
	return requestedPlan, nil
}

// UserPlans returns all user's plans by user id.
func (s *PlanStorage) UserPlans(
	_ context.Context,
	userID uuid.UUID,
) ([]*plan.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*plan.Plan
	for _, onePlan := range s.data.GetAll() {
		if onePlan.UserID() == userID {
			result = append(result, onePlan)
		}
	}
	if len(result) == 0 {
		return nil, plan.ErrNoPlanFound
	}
	return result, nil
}
