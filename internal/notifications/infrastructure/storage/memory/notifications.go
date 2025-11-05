package memory

import (
	"context"
	"sync"

	notifications "github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/notifications"
	"github.com/google/uuid"
)

// NotificationsStorage is a storage for notifications.
type NotificationsStorage struct {
	data  *Cache[*notifications.Notification]
	count uint

	mu *sync.RWMutex
}

// NewNotificationsStorage returns a new NotificationsStorage.
func NewNotificationsStorage() *NotificationsStorage {
	return &NotificationsStorage{
		data:  NewCache[*notifications.Notification](),
		count: 0,
		mu:    &sync.RWMutex{},
	}
}

// Create creates a new medication in memory.
func (s *NotificationsStorage) Create(
	_ context.Context,
	notification *notifications.Notification,
) (*notifications.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.count++
	s.data.Set(notification.GetID().String(), notification)
	return notification, nil
}

// GetByID returns a medication by id.
func (s *NotificationsStorage) GetByID(
	_ context.Context,
	notificationID uuid.UUID,
) (*notifications.Notification, error) {
	drug, ok := s.data.Get(notificationID.String())
	if !ok {
		return nil, notifications.ErrNoNotificationsFound
	}
	return drug, nil
}
