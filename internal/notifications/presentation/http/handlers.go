// Package http is a package for http handlers
package http

import (
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
			IsActive:  serviceResponse.IsActive,
		},
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// DeleteSubscription delete a subscription one time for every device.
func (h *NotificationsHandlers) DeleteSubscriptionGin(c *gin.Context) {
	if _, err := httputil.GetAuthFromCtx(c.Request); err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	slugUserID := c.Param(SlugID)
	if slugUserID == "" {
		h.logger.Error("Subscription ID not found in path params")
		c.JSON(http.StatusBadRequest, api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgMissingSlug,
			Body:       struct{}{},
		})
		return
	}

	serviceRequest := &application.DeleteSubscriptionCommand{
		UserID: slugUserID,
	}
	_, err := h.app.DeleteSubscription.Execute(c.Request.Context(), serviceRequest)
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
		SendAt: reqJSON.SendAt,
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
