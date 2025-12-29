// Package notification is a package for adapting application logic for sending notifications.
package notification

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/notification"
	client "github.com/FSO-VK/final-project-vk-backend/internal/utils/notification_client"
)

type NotificationProvider struct {
	client *client.NotificationClient
}

func NewNotificationProvider(c *client.NotificationClient) *NotificationProvider {
	return &NotificationProvider{client: c}
}

func (a *NotificationProvider) SendNotification(
	ctx context.Context,
	info notification.NotificationInfo,
) error {
	return a.client.SendNotification(ctx, client.NotificationInfo{
		UserID: info.UserID,
		Title:  info.Title,
		Body:   info.Body,
	})
}
