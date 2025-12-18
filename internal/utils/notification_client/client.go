// Package notificationclient implements NotificationService interface for sending notification to Notification service.
package notificationclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// NotificationClient implements NotificationService.
type NotificationClient struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

// NewNotificationClient creates a new NotificationClient.
func NewNotificationClient(cfg ClientConfig, logger *logrus.Entry) *NotificationClient {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &NotificationClient{client: client, cfg: cfg, logger: logger}
}

type NotificationExpectedResponse struct {
	StatusCode int  `json:"statusCode"`
	Body       Body `json:"body"`
}

type Body struct{}

type NotificationInfo struct {
	UserID uuid.UUID `json:"userId"`
	Title  string    `json:"title"`
	Body   string    `json:"body"`
}

// SendNotification implements NotificationService interface and sends a notification.
func (h *NotificationClient) SendNotification(
	ctx context.Context,
	info NotificationInfo,
) error {
	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}

	jsonBody, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal notification body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		h.cfg.Method,
		h.cfg.Endpoint,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		h.logger.WithError(err).Warn("notification API request failed")
		return ErrNotificationServiceUnavailable
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		h.logger.Warnf("notification service responded with %d", resp.StatusCode)
		return ErrBadResponse
	}

	var parsed NotificationExpectedResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		h.logger.WithError(err).Error("failed to decode notification API response")
		return fmt.Errorf("%w: %w", ErrBadResponse, err)
	}

	if parsed.StatusCode != http.StatusOK {
		return ErrBadResponse
	}

	return nil
}
