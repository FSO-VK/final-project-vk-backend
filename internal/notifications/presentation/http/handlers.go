// Package http is a package for http handlers
package http

import (
	"fmt"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
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
type GetVapidPublicKeyJSONRequest struct {
}

// GetVapidPublicKeyJSONResponse is a response for CreateSubscription.
type GetVapidPublicKeyJSONResponse struct {
	VapidPublicKey string `json:"vapidPublicKey"`
}

// GetVapidPublicKey adds a medication.
func (h *NotificationsHandlers) GetVapidPublicKey(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	fmt.Println(auth)
}

// CreateSubscriptionJSONRequest is a request for CreateSubscription.
type  CreateSubscriptionJSONResponse struct {
	SubscriptionObject
}

// CreateSubscription create a subscription one time for every device.
func (h *NotificationsHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	fmt.Println(auth)
}

// DeleteSubscriptionJSONRequest is a request for DeleteSubscription.
type DeleteSubscriptionJSONRequest struct {
}

// DeleteSubscriptionJSONResponse is a response for DeleteSubscription.
type DeleteSubscriptionJSONResponse struct {
}

// DeleteSubscription delete a subscription one time for every device.
func (h *NotificationsHandlers) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	fmt.Println(auth)
}

// SendNotificationJSONRequest is a request for SendNotification.
type SendNotificationJSONRequest struct {
	// TODO
}

// SendNotificationJSONResponse is a response for SendNotification.
type SendNotificationJSONResponse struct {
}

// SendNotification delete a subscription one time for every device.
func (h *NotificationsHandlers) SendNotification(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	fmt.Println(auth)
}

// InteractWithNotificationJSONRequest is a request for InteractWithNotification.
type InteractWithNotificationJSONRequest struct {
	Action string `json:"action"` // snooze, dismiss, apply
}

// InteractWithNotificationJSONResponse is a response for InteractWithNotification.
type InteractWithNotificationJSONResponse struct {
}

// InteractWithNotification delete a subscription one time for every device.
func (h *NotificationsHandlers) InteractWithNotification(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	fmt.Println(auth)
}
