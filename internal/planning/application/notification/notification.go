package notification

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	SendNotification(ctx context.Context, notificationInfo NotificationInfo) error
}

type NotificationInfo struct {
	UserID uuid.UUID
	Title  string
	Body   string
}
