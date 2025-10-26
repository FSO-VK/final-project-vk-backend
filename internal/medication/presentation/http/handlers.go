// Package http is a package for http handlers
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// SlugID is a slug for id.
	SlugID = "id"
)

// MedicationHandlers is a handler for Medication.
type MedicationHandlers struct {
	app    *application.MedicationApplication
	logger *logrus.Entry
}

// NewHandlers creates a new MedicationHandlers.
func NewHandlers(
	app *application.MedicationApplication,
	logger *logrus.Entry,
) *MedicationHandlers {
	return &MedicationHandlers{
		app:    app,
		logger: logger,
	}
}

// AddMedicationJSONRequest is a request for AddMedication.
type AddMedicationJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// AddMedicationJSONResponse is a response for AddMedication.
type AddMedicationJSONResponse struct {
	ID string `json:"id"`
}

// AddMedication adds a medication.
func (h *MedicationHandlers) AddMedication(w http.ResponseWriter, r *http.Request) {
	var reqJSON *AddMedicationJSONRequest

	var body bytes.Buffer
	_, err := body.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgFailedToReadBody,
			Body:       nil,
		})

		return
	}

	err = json.Unmarshal(body.Bytes(), &reqJSON)
	if err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgFailedToUnmarshal,
			Body:       nil,
		})
		return
	}

	serviceRequest := &application.AddMedicationCommand{
		Name:         reqJSON.Name,
		CategoriesID: nil,
		Items:        reqJSON.Items,
		ItemsUnit:    reqJSON.ItemsUnit,
		Expires:      reqJSON.Expiration,
	}

	serviceResponse, err := h.app.AddMedication.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add medication")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       nil,
			Error:      MsgFailedToAddMedication,
		})

		return
	}
	response := &AddMedicationJSONResponse{
		ID: strconv.FormatUint(uint64(serviceResponse.ID), 10),
	}
	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// UpdateMedicationJSONRequest is a request for UpdateMedication.
type UpdateMedicationJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// UpdateMedicationJSONResponse is a response for UpdateMedication.
type UpdateMedicationJSONResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// UpdateMedication updates a medication.
func (h *MedicationHandlers) UpdateMedication(w http.ResponseWriter, r *http.Request) {
	var reqJSON *UpdateMedicationJSONRequest

	vars := mux.Vars(r)
	id := vars[SlugID]
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to parse id")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailToParseID,
		})
		return
	}

	var body bytes.Buffer
	_, err = body.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailedToReadBody,
		})
		return
	}

	err = json.Unmarshal(body.Bytes(), &reqJSON)
	if err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailedToUnmarshal,
		})
		return
	}

	serviceRequest := &application.UpdateMedicationCommand{
		ID:           uint(idUint),
		Name:         reqJSON.Name,
		CategoriesID: nil,
		Items:        reqJSON.Items,
		ItemsUnit:    reqJSON.ItemsUnit,
		Expires:      reqJSON.Expiration,
	}

	serviceResponse, err := h.app.UpdateMedication.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update medication")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToUpdateMedication,
		})

		return
	}

	response := &UpdateMedicationJSONResponse{
		ID:         strconv.FormatUint(uint64(serviceResponse.ID), 10),
		Name:       serviceResponse.Name,
		Items:      serviceResponse.Items,
		ItemsUnit:  serviceResponse.ItemsUnit,
		Expiration: serviceResponse.Expires,
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// DeleteMedication deletes a medication.
func (h *MedicationHandlers) DeleteMedication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[SlugID]

	serviceRequest := &application.DeleteMedicationCommand{
		ID: id,
	}

	_, err := h.app.DeleteMedication.Execute(
		r.Context(),
		serviceRequest,
	)

	if err != nil {
		if errors.Is(application.ErrDeleteInvalidUuidFormat, err) {
			h.logger.WithError(err).Error("Failed to parse")
			w.WriteHeader(http.StatusBadRequest)
			_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
				StatusCode: http.StatusBadRequest,
				Body:       struct{}{},
				Error:      MsgFailToParseID,
			})
		} else {
			h.logger.WithError(err).Error("Failed to delete medication")
			w.WriteHeader(http.StatusInternalServerError)

			_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
				StatusCode: http.StatusInternalServerError,
				Body:       struct{}{},
				Error:      MsgFailedToDeleteMedication,
			})
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[struct{}]{
		StatusCode: http.StatusOK,
		Body:       struct{}{},
		Error:      "",
	})
}

// GetMedicationListItem returns a list of medications.
type GetMedicationListItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// GetMedicationListJSONResponse returns a list of medications.
type GetMedicationListJSONResponse struct {
	MedicationList []GetMedicationListItem `json:"medicationList"`
}

// GetMedicationList returns a list of medications.
func (h *MedicationHandlers) GetMedicationList(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := h.app.GetMedicationList.Execute(
		r.Context(),
		&application.GetMedicationListCommand{},
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get medication list")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedicationList,
		})
		return
	}

	response := &GetMedicationListJSONResponse{
		MedicationList: make([]GetMedicationListItem, 0),
	}

	for _, item := range serviceResponse.List {
		response.MedicationList = append(response.MedicationList, GetMedicationListItem{
			ID:         strconv.FormatUint(uint64(item.ID), 10),
			Name:       item.Name,
			Items:      item.Items,
			ItemsUnit:  item.ItemsUnit,
			Expiration: item.Expires,
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}
