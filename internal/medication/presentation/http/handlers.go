// Package http is a package for http handlers
package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// SlugID is a slug for id.
	SlugID = "id"
)

// MedicineHandlers is a handler for Medicine.
type MedicineHandlers struct {
	app    *application.MedicineApplication
	logger *logrus.Entry
}

// NewHandlers creates a new MedicineHandlers.
func NewHandlers(
	app *application.MedicineApplication,
	logger *logrus.Entry,
) *MedicineHandlers {
	return &MedicineHandlers{
		app:    app,
		logger: logger,
	}
}

// AddMedicineJSONRequest is a request for AddMedicine.
type AddMedicineJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// AddMedicine adds a medicine.
func (h *MedicineHandlers) AddMedicine(w http.ResponseWriter, r *http.Request) {
	var reqJSON *AddMedicineJSONRequest

	var body bytes.Buffer
	_, err := body.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
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

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgFailedToUnmarshal,
			Body:       nil,
		})
		return
	}

	serviceRequest := &application.AddMedicineCommand{
		Name:         reqJSON.Name,
		CategoriesID: nil,
		Items:        reqJSON.Items,
		ItemsUnit:    reqJSON.ItemsUnit,
		Expires:      reqJSON.Expiration,
	}

	serviceResponse, err := h.app.AddMedicine.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       nil,
			Error:      MsgFailedToAddMedicine,
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       serviceResponse,
		Error:      "",
	})
}

// UpdateMedicineJSONRequest is a request for UpdateMedicine.
type UpdateMedicineJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// UpdateMedicineJSONResponse is a response for UpdateMedicine.
type UpdateMedicineJSONResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// UpdateMedicine updates a medicine.
func (h *MedicineHandlers) UpdateMedicine(w http.ResponseWriter, r *http.Request) {
	var reqJSON *UpdateMedicineJSONRequest

	vars := mux.Vars(r)
	id := vars[SlugID]
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to parse id")
		w.WriteHeader(http.StatusBadRequest)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
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

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
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

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailedToUnmarshal,
		})
		return
	}

	serviceRequest := &application.UpdateMedicineCommand{
		ID:           uint(idUint),
		Name:         reqJSON.Name,
		CategoriesID: nil,
		Items:        reqJSON.Items,
		ItemsUnit:    reqJSON.ItemsUnit,
		Expires:      reqJSON.Expiration,
	}

	serviceResponse, err := h.app.UpdateMedicine.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToUpdateMedicine,
		})

		return
	}

	response := &UpdateMedicineJSONResponse{
		ID:         strconv.FormatUint(uint64(serviceResponse.ID), 10),
		Name:       serviceResponse.Name,
		Items:      serviceResponse.Items,
		ItemsUnit:  serviceResponse.ItemsUnit,
		Expiration: serviceResponse.Expires,
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// DeleteMedicine deletes a medicine.
func (h *MedicineHandlers) DeleteMedicine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[SlugID]
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to parse id")
		w.WriteHeader(http.StatusBadRequest)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailToParseID,
		})
		return
	}

	serviceRequest := &application.DeleteMedicineCommand{
		ID: uint(idUint),
	}

	_, err = h.app.DeleteMedicine.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToDeleteMedicine,
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[struct{}]{
		StatusCode: http.StatusOK,
		Body:       struct{}{},
		Error:      "",
	})
}

// GetMedicineListItem returns a list of medicines.
type GetMedicineListItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

// GetMedicineListJSONResponse returns a list of medicines.
type GetMedicineListJSONResponse struct {
	MedicineList []GetMedicineListItem `json:"medicineList"`
}

// GetMedicineList returns a list of medicines.
func (h *MedicineHandlers) GetMedicineList(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := h.app.GetMedicineList.Execute(
		r.Context(),
		&application.GetMedicineListCommand{},
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get medicine list")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedicineList,
		})
		return
	}

	response := &GetMedicineListJSONResponse{
		MedicineList: make([]GetMedicineListItem, 0),
	}

	for _, item := range serviceResponse.List {
		response.MedicineList = append(response.MedicineList, GetMedicineListItem{
			ID:         strconv.FormatUint(uint64(item.ID), 10),
			Name:       item.Name,
			Items:      item.Items,
			ItemsUnit:  item.ItemsUnit,
			Expiration: item.Expires,
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}
