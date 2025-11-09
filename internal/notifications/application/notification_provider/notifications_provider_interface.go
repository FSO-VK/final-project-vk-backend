// Package notificationprovider is a package for interface for notifications client.
package notificationprovider

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/google/uuid"
)

// PushNotificationClient is an interface for notifications client.
type NotificationProvider interface {
	PushNotification(
		pushInfo *Notification,
		subscriptionInfo *subscriptions.PushSubscription,
	) error
}

// Notification represents a push notification to be sent.
type Notification struct {
	ID             uuid.UUID
	SubscriptionID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Body           string
}

// NewNotification creates a new notification.
func NewNotification(
	id uuid.UUID,
	subscriptionID uuid.UUID,
	userID uuid.UUID,
	title string,
	body string,
) *Notification {
	return &Notification{
		ID:             id,
		SubscriptionID: subscriptionID,
		UserID:         userID,
		Title:          title,
		Body:           body,
	}
}
