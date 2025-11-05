// Package notifications is a domain layer for notifications.
package notifications

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a push notification to be sent.
type Notification struct {
	id             uuid.UUID
	subscriptionID uuid.UUID
	userID         uuid.UUID
	title          string
	body           string
	sendAt         time.Time
	planningID     uuid.UUID
}

// NewNotification creates a new notification.
func NewNotification(
	id uuid.UUID,
	subscriptionID uuid.UUID,
	userID uuid.UUID,
	title string,
	body string,
	sendAt time.Time,
	planningID uuid.UUID,
) *Notification {
	return &Notification{
		id:             id,
		subscriptionID: subscriptionID,
		userID:         userID,
		title:          title,
		body:           body,
		sendAt:         sendAt,
		planningID:     planningID,
	}
}

// GetID returns the unique identifier of the notification.
func (n *Notification) GetID() uuid.UUID {
	return n.id
}

// SetID sets the unique identifier of the notification.
func (n *Notification) SetID(id uuid.UUID) {
	n.id = id
}

// GetSubscriptionID returns the subscription ID.
func (n *Notification) GetSubscriptionID() uuid.UUID {
	return n.subscriptionID
}

// SetSubscriptionID sets the subscription ID.
func (n *Notification) SetSubscriptionID(subscriptionID uuid.UUID) {
	n.subscriptionID = subscriptionID
}

// GetUserID returns the user identifier.
func (n *Notification) GetUserID() uuid.UUID {
	return n.userID
}

// SetUserID sets the user identifier.
func (n *Notification) SetUserID(userID uuid.UUID) {
	n.userID = userID
}

// GetTitle returns the notification title.
func (n *Notification) GetTitle() string {
	return n.title
}

// SetTitle sets the notification title.
func (n *Notification) SetTitle(title string) {
	n.title = title
}

// GetBody returns the notification body.
func (n *Notification) GetBody() string {
	return n.body
}

// SetBody sets the notification body.
func (n *Notification) SetBody(body string) {
	n.body = body
}

// GetSendAt returns the send time.
func (n *Notification) GetSendAt() time.Time {
	return n.sendAt
}

// SetSendAt sets the send time.
func (n *Notification) SetSendAt(sendAt time.Time) {
	n.sendAt = sendAt
}

// GetPlanningID returns the planning ID.
func (n *Notification) GetPlanningID() uuid.UUID {
	return n.planningID
}

// SetPlanningID sets the planning ID.
func (n *Notification) SetPlanningID(planningID uuid.UUID) {
	n.planningID = planningID
}
