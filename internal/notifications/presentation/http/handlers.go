// Package http is a package for http handlers
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// SlugID is a slug for id.
	SlugID = "id"
)

// will be removed after 06.12.
var (
	errNotificationMarshal   = errors.New("failed to marshal notification")
	errNotificationSend      = errors.New("failed to send notification")
	errNotificationBadStatus = errors.New("notification service returned bad status")
)

// NotificationsHandlers is a handler for Notifications.
type NotificationsHandlers struct {
	app    *application.NotificationsApplication
	logger *logrus.Entry
}

// NewHandlers creates a new NotificationsHandlers.
func NewHandlers(
	app *application.NotificationsApplication,
	logger *logrus.Entry,
) *NotificationsHandlers {
	return &NotificationsHandlers{
		app:    app,
		logger: logger,
	}
}

// GetVapidPublicKeyJSONRequest is a request for CreateSubscription.
type GetVapidPublicKeyJSONRequest struct{}

// GetVapidPublicKeyJSONResponse is a response for CreateSubscription.
type GetVapidPublicKeyJSONResponse struct {
	VapidPublicKey string `json:"vapidPublicKey"`
}

// GetVapidPublicKey adds a medication.
func (h *NotificationsHandlers) GetVapidPublicKeyGin(c *gin.Context) {
	if _, err := httputil.GetAuthFromCtx(c.Request); err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	serviceRequest := &application.GetVapidPublicKeyCommand{}
	serviceResponse, err := h.app.GetVapidPublicKey.Execute(c.Request.Context(), serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get vapid public key")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetVapidPublicKey,
		})
		return
	}
	response := &GetVapidPublicKeyJSONResponse{
		VapidPublicKey: serviceResponse.PublicKey,
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// CreateSubscriptionJSONRequest is a request for CreateSubscription.
type CreateSubscriptionJSONRequest struct {
	// embedded struct
	SubscriptionObject `json:",inline"`
}

// CreateSubscriptionJSONResponse is a request for CreateSubscription.
type CreateSubscriptionJSONResponse struct {
	// embedded struct
	PushSubscriptionInfo `json:",inline"`
}

// CreateSubscription create a subscription one time for every device.
func (h *NotificationsHandlers) CreateSubscriptionGin(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	var reqJSON CreateSubscriptionJSONRequest
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		h.logger.WithError(err).Error("Failed to bind request body")
		c.JSON(http.StatusBadRequest, api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		})
		return
	}

	command := &application.CreateSubscriptionCommand{
		UserID: auth.UserID,
		SendInfo: application.SendInfo{
			Endpoint: reqJSON.SendInfo.Endpoint,
			Keys: application.Keys{
				P256dh: reqJSON.SendInfo.Keys.P256dh,
				Auth:   reqJSON.SendInfo.Keys.Auth,
			},
		},
		UserAgent: reqJSON.UserAgent,
	}

	serviceResponse, err := h.app.CreateSubscription.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create subscription")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToCreateSubscription,
		})
		return
	}

	response := &CreateSubscriptionJSONResponse{
		PushSubscriptionInfo: PushSubscriptionInfo{
			ID:     serviceResponse.ID,
			UserID: serviceResponse.UserID,
			SendInfo: SendInfo{
				Endpoint: serviceResponse.SendInfo.Endpoint,
				Keys: Keys{
					P256dh: serviceResponse.SendInfo.Keys.P256dh,
					Auth:   serviceResponse.SendInfo.Keys.Auth,
				},
			},
			UserAgent: serviceResponse.UserAgent,
		},
	}
	go func(ctx context.Context) {
		if err := sendWelcomeNotification(ctx, auth.UserID); err != nil {
			h.logger.WithError(err).Warn("Failed to send welcome notification")
		}
	}(c.Request.Context())

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

type WelcomeNotificationRequest struct {
	UserID string `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func sendWelcomeNotification(ctx context.Context, userID string) error {
	notification := WelcomeNotificationRequest{
		UserID: userID,
		Title:  "Приветственное напоминание",
		Body:   "Спасибо, что включили уведомления, рады видеть вас в нашем сервисе!",
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("%w: %w", errNotificationMarshal, err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://myhealthbox.ddns.net/api/v1/notification/send",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("%w: %w", errNotificationSend, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", errNotificationSend, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", errNotificationBadStatus, resp.StatusCode)
	}

	return nil
}

// DeleteSubscription delete a subscription one time for every device.
func (h *NotificationsHandlers) DeleteSubscriptionGin(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	serviceRequest := &application.DeleteSubscriptionCommand{
		UserID: auth.UserID,
	}
	_, err = h.app.DeleteSubscription.Execute(c.Request.Context(), serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete subscription")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToDeleteSubscription,
		})
		return
	}

	c.JSON(http.StatusOK, api.Response[struct{}]{
		StatusCode: http.StatusOK,
		Body:       struct{}{},
		Error:      "",
	})
}

// SendNotificationJSONRequest is a request for SendNotification.
type SendNotificationJSONRequest struct {
	// embedded struct
	PushNotificationObject `json:",inline"`

	UserID string `json:"userId"`
}

// SendNotificationJSONResponse is a response for SendNotification.
type SendNotificationJSONResponse struct{}

// SendNotification delete a subscription one time for every device.
func (h *NotificationsHandlers) SendNotificationGin(c *gin.Context) {
	var reqJSON SendNotificationJSONRequest
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		h.logger.WithError(err).Error("Failed to bind request body")
		c.JSON(http.StatusBadRequest, api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		})
		return
	}

	command := &application.SendNotificationCommand{
		UserID: reqJSON.UserID,
		Title:  reqJSON.Title,
		Body:   reqJSON.Body,
	}

	_, err := h.app.SendNotification.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send notification")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToSendNotification,
		})
		return
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       struct{}{},
		Error:      "",
	})
}
