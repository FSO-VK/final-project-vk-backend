// Package notificationprovider is a package for interface for notifications client.
package notificationprovider

import (
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/google/uuid"
)

// PushNotificationClient is an interface for data matrix client.
type NotificationProvider interface {
	PushNotification(
		pushInfo *PushNotification,
		subscriptionInfo *subscriptions.PushSubscription,
	) error
}

// PushNotification represents a push notification to be sent.
type PushNotification struct {
	id             uuid.UUID
	subscriptionID uuid.UUID
	userID         uuid.UUID
	title          string
	body           string
	sendAt         time.Time
}

// NewNotification creates a new notification.
func NewNotification(
	id uuid.UUID,
	subscriptionID uuid.UUID,
	userID uuid.UUID,
	title string,
	body string,
	sendAt time.Time,
) *PushNotification {
	return &PushNotification{
		id:             id,
		subscriptionID: subscriptionID,
		userID:         userID,
		title:          title,
		body:           body,
		sendAt:         sendAt,
	}
}

// GetID returns the unique identifier of the notification.
func (n *PushNotification) GetID() uuid.UUID {
	return n.id
}

// SetID sets the unique identifier of the notification.
func (n *PushNotification) SetID(id uuid.UUID) {
	n.id = id
}

// GetSubscriptionID returns the subscription ID.
func (n *PushNotification) GetSubscriptionID() uuid.UUID {
	return n.subscriptionID
}

// SetSubscriptionID sets the subscription ID.
func (n *PushNotification) SetSubscriptionID(subscriptionID uuid.UUID) {
	n.subscriptionID = subscriptionID
}

// GetUserID returns the user identifier.
func (n *PushNotification) GetUserID() uuid.UUID {
	return n.userID
}

// SetUserID sets the user identifier.
func (n *PushNotification) SetUserID(userID uuid.UUID) {
	n.userID = userID
}

// GetTitle returns the notification title.
func (n *PushNotification) GetTitle() string {
	return n.title
}

// SetTitle sets the notification title.
func (n *PushNotification) SetTitle(title string) {
	n.title = title
}

// GetBody returns the notification body.
func (n *PushNotification) GetBody() string {
	return n.body
}

// SetBody sets the notification body.
func (n *PushNotification) SetBody(body string) {
	n.body = body
}

// GetSendAt returns the send time.
func (n *PushNotification) GetSendAt() time.Time {
	return n.sendAt
}

// SetSendAt sets the send time.
func (n *PushNotification) SetSendAt(sendAt time.Time) {
	n.sendAt = sendAt
}
