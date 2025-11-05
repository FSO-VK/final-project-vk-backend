package memory

import (
	"context"
	"sync"

	subscriptions "github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/google/uuid"
)

// SubscriptionsStorage is a storage for subscriptions.
type SubscriptionsStorage struct {
	data  *Cache[*subscriptions.Subscription]
	count uint

	mu *sync.RWMutex
}

// NewSubscriptionsStorage returns a new SubscriptionsStorage.
func NewSubscriptionsStorage() *SubscriptionsStorage {
	return &SubscriptionsStorage{
		data:  NewCache[*subscriptions.Subscription](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Create creates a new medication in memory.
func (s *SubscriptionsStorage) Create(
	_ context.Context,
	subscription *subscriptions.Subscription,
) (*subscriptions.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.count++
	s.data.Set(subscription.GetID().String(), subscription)
	return subscription, nil
}

// GetByID returns a medication by id.
func (s *SubscriptionsStorage) GetByID(
	_ context.Context,
	subscriptionID uuid.UUID,
) (*subscriptions.Subscription, error) {
	drug, ok := s.data.Get(subscriptionID.String())
	if !ok {
		return nil, subscriptions.ErrNoSubscriptionsFound
	}
	return drug, nil
}
