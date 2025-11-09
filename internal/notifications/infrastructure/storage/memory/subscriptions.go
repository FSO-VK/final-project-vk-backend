package memory

import (
	"context"
	"sync"

	subscriptions "github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

// SubscriptionsStorage is a storage for subscriptions.
type SubscriptionsStorage struct {
	data  *cache.Cache[*subscriptions.PushSubscription]
	count uint

	mu *sync.RWMutex
}

// NewSubscriptionsStorage returns a new SubscriptionsStorage.
func NewSubscriptionsStorage() *SubscriptionsStorage {
	return &SubscriptionsStorage{
		data:  cache.NewCache[*subscriptions.PushSubscription](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// CreateSubscription creates a new subscription in memory or updates it.
func (s *SubscriptionsStorage) CreateSubscription(
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

// GetSubscriptionsByUserID returns all subscriptions with the same user id.
func (s *SubscriptionsStorage) GetSubscriptionsByUserID(
	_ context.Context,
	userID uuid.UUID,
) ([]*subscriptions.PushSubscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*subscriptions.PushSubscription
	for _, subscription := range s.data.GetAll() {
		if subscription.GetUserID() == userID {
			result = append(result, subscription)
		}
	}
	if len(result) == 0 {
		return nil, subscriptions.ErrNoSubscriptionsFound
	}
	return result, nil
}

// DeleteSubscription removes a subscription from memory.
func (s *SubscriptionsStorage) DeleteSubscription(
	_ context.Context,
	subscriptionID uuid.UUID,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := subscriptionID.String()
	if _, exists := s.data.Get(id); exists {
		s.data.Delete(id)
		s.count--
		return nil
	}
	return subscriptions.ErrNoSubscriptionsFound
}
