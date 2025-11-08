// Package client implements NotificationProvider interface for sending pushing notifications.
package client

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
	pushInfo *notificationProvider.PushNotification,
	subscriptionInfo *subscriptions.PushSubscription,
) error {
	if pushInfo == nil || subscriptionInfo == nil {
		return ErrBadRequest
	}
	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
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
		"title": pushInfo.GetTitle(),
		"body":  pushInfo.GetBody(),
		"data": map[string]interface{}{
			"userID": pushInfo.GetUserID().String(),
		},
	}

	payload, err := json.Marshal(notificationPayload)
	if err != nil {
		h.logger.WithError(err).Error("failed to marshal notification payload")
		return err
	}

	options := &webpush.Options{
		HTTPClient:      h.client,
		Subscriber:      "https://myhealthbox.ddns.net/",
		TTL:             600,
		Urgency:         webpush.UrgencyNormal, // normal priority would not receive notification when low battery
		VAPIDPublicKey:  h.cfg.VapidPublicKey,
		VAPIDPrivateKey: h.cfg.VapidPrivateKey,
	}

	resp, err := webpush.SendNotificationWithContext(ctx, payload, webpushSubscription, options)
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
			"subscription_id": pushInfo.GetSubscriptionID(),
			"user_id":         pushInfo.GetUserID(),
			"status":          resp.Status,
		}).Info("push notification sent successfully")
		return nil
	}
	h.logger.WithFields(logrus.Fields{
		"subscription_id": pushInfo.GetSubscriptionID(),
		"status_code":     resp.StatusCode,
		"status":          resp.Status,
	}).Warn("push notification failed to sent with non-200 status")

	return ErrPushServiceUnavailable
}
