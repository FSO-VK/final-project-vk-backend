// Package http is a package for http handlers
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gorilla/mux"
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

// AddNotificationsJSONRequest is a request for AddNotifications.
type AddMedicationJSONRequest struct {
	BodyCommonObject `json:",inline"`
}

// AddMedicationJSONResponse is a response for AddMedication.
type AddMedicationJSONResponse struct {
	// embedded struct
	BodyCommonObject `json:",inline"`

	ID string `json:"id"`
}

// AddMedication adds a medication.
func (h *MedicationHandlers) AddNotifications(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	
}


// DeleteMedication deletes a medication.
func (h *MedicationHandlers) DeleteMedication(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}
	
}

// GetMedicationBoxItem returns a Box of medications.
type GetMedicationBoxItem struct {
	// embedded struct
	BodyCommonObject `json:",inline"`

	ID string `json:"id"`
}

// GetMedicationBoxJSONResponse returns a Box of medications.
type GetMedicationBoxJSONResponse struct {
	MedicationBox []GetMedicationBoxItem `json:"medicationBox"`
}

// GetMedicationBox returns a Box of medications.
func (h *MedicationHandlers) GetMedicationBox(w http.ResponseWriter, r *http.Request) {
	authorization, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}
	
}
