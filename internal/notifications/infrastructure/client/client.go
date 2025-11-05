// Package client implements PushNotificationClient interface for getting pushing notifications.
package client

import (
	"context"
	"net/http"

	clientInterface "github.com/FSO-VK/final-project-vk-backend/internal/notifications/application/notifications_client"
	"github.com/sirupsen/logrus"
)

// PushClientAPI implements PushNotificationClient.
type PushClientAPI struct {
	client *http.Client
	cfg    PushClient
	logger *logrus.Entry
}

// NewNewPushClient creates a new NewPushClient.
func NewPushClientAPI(cfg PushClient, logger *logrus.Entry) *PushClientAPI {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &PushClientAPI{client: client, cfg: cfg, logger: logger}
}

// PushNotification implements PushNotificationClient interface.
func (h *PushClientAPI) PushNotification(
	pushInfo *clientInterface.PushNotificationInfo,
) (*clientInterface.NotificationResponse, error) {
	if pushInfo == nil {
		return nil, ErrBadRequest
	}
	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}
	// TODO

	return nil, nil
}
