package notification

import (
	"context"

	"github.com/google/uuid"
)

// NotificationService is service for sending notifications.
type NotificationService interface {
	SendNotification(ctx context.Context, notificationInfo NotificationInfo) error
}

// NotificationInfo contains information of notification.
type NotificationInfo struct {
	UserID uuid.UUID
	Title  string
	Body   string
}
