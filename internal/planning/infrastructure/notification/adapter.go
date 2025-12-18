// Package notification is a package for adapting application logic for sending notifications.
package notification

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/notification"
	client "github.com/FSO-VK/final-project-vk-backend/internal/utils/notification_client"
)

type Adapter struct {
	client *client.NotificationClient
}

func NewAdapter(c *client.NotificationClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) SendNotification(
	ctx context.Context,
	info notification.NotificationInfo,
) error {
	return a.client.SendNotification(ctx, client.NotificationInfo{
		UserID: info.UserID,
		Title:  info.Title,
		Body:   info.Body,
	})
}
