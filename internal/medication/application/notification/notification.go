package notification

import (
	"context"

	"github.com/google/uuid"
)

// NotificationService is an interface for sending notification to Notification service.
type NotificationService interface {
	SendNotification(ctx context.Context, notificationInfo NotificationInfo) error
}

// NotificationInfo is a struct for sending notification.
type NotificationInfo struct {
	UserID uuid.UUID
	Title  string
	Body   string
}
