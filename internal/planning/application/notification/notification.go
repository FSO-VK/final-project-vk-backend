package notification

import "github.com/google/uuid"

type NotificationService interface {
	SendNotification(notificationInfo NotificationInfo) error
}

type NotificationInfo struct {
	UserID uuid.UUID
	Title  string
	Body   string
}
