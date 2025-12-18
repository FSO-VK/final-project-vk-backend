// Package psuhprovider implements NotificationProvider interface for sending pushing notifications.
package psuhprovider

import (
	"context"
	"encoding/json"
	"net/http"

	notificationProvider "github.com/FSO-VK/final-project-vk-backend/internal/notifications/application/notification_provider"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/sirupsen/logrus"
)

// PushNotificationProvider implements NotificationProvider.
type PushNotificationProvider struct {
	client *http.Client
	cfg    PushClient
	logger *logrus.Entry
}

// NewPushNotificationProvider creates a new NewPushClient.
func NewPushNotificationProvider(cfg PushClient, logger *logrus.Entry) *PushNotificationProvider {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &PushNotificationProvider{client: client, cfg: cfg, logger: logger}
}

// PushNotification implements NotificationProvider interface.
func (h *PushNotificationProvider) PushNotification(
	pushInfo *notificationProvider.Notification,
	subscriptionInfo *subscriptions.PushSubscription,
) error {
	if pushInfo == nil || subscriptionInfo == nil {
		return ErrBadRequest
	}
	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		_, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}

	webpushSubscription := &webpush.Subscription{
		Endpoint: subscriptionInfo.GetEndpoint(),
		Keys: webpush.Keys{
			P256dh: subscriptionInfo.GetP256dh(),
			Auth:   subscriptionInfo.GetAuth(),
		},
	}

	notificationPayload := map[string]interface{}{
		"title": pushInfo.Title,
		"body":  pushInfo.Body,
		"data": map[string]interface{}{
			"userID": pushInfo.UserID.String(),
		},
	}

	payload, err := json.Marshal(notificationPayload)
	if err != nil {
		h.logger.WithError(err).Error("failed to marshal notification payload")
		return err
	}

	resp, err := webpush.SendNotification(payload, webpushSubscription, &webpush.Options{
		Subscriber:      h.cfg.Subscriber,
		VAPIDPublicKey:  h.cfg.VapidPublicKey,
		VAPIDPrivateKey: h.cfg.VapidPrivateKey,
		TTL:             int(h.cfg.Timeout),
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to send push notification")
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.WithError(err).Debug("failed to close response body")
		}
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		h.logger.WithFields(logrus.Fields{
			"subscription_id": pushInfo.ID,
			"user_id":         pushInfo.UserID,
			"status":          resp.Status,
		}).Info("push notification sent successfully")
		return nil
	}
	h.logger.WithFields(logrus.Fields{
		"subscription_id": pushInfo.ID,
		"status_code":     resp.StatusCode,
		"status":          resp.Status,
	}).Warn("push notification failed to sent with non-200 status")

	return ErrPushServiceUnavailable
}
