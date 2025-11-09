package memory

import (
	"context"
	"sync"

	subscriptions "github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/google/uuid"
)

// SubscriptionsStorage is a storage for subscriptions.
type SubscriptionsStorage struct {
	data  *Cache[*subscriptions.PushSubscription]
	count uint

	mu *sync.RWMutex
}

// NewSubscriptionsStorage returns a new SubscriptionsStorage.
func NewSubscriptionsStorage() *SubscriptionsStorage {
	return &SubscriptionsStorage{
		data:  NewCache[*subscriptions.PushSubscription](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// SetSubscription creates a new subscription in memory or updates it.
func (s *SubscriptionsStorage) SetSubscription(
	_ context.Context,
	subscription *subscriptions.PushSubscription,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	subscriptionID := subscription.GetID().String()
	if _, exists := s.data.Get(subscriptionID); exists {
		s.data.Set(subscriptionID, subscription)
		return nil
	}

	s.count++
	s.data.Set(subscriptionID, subscription)
	return nil
}

// GetSubscriptionByID returns a subscription by id.
func (s *SubscriptionsStorage) GetSubscriptionByID(
	_ context.Context,
	subscriptionID uuid.UUID,
) (*subscriptions.PushSubscription, error) {
	drug, ok := s.data.Get(subscriptionID.String())
	if !ok {
		return nil, subscriptions.ErrNoSubscriptionsFound
	}
	return drug, nil
}
