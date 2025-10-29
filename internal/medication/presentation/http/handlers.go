// Package http is a package for http handlers
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

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
	BodyCommonObject `json:",inline"`
}

// AddMedicationJSONResponse is a response for AddMedication.
type AddMedicationJSONResponse struct {
	// embedded struct
	BodyCommonObject `json:",inline"`

	ID string `json:"id"`
}

// AddMedication adds a medication.
func (h *MedicationHandlers) AddMedication(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       nil,
		})
		return
	}

	var reqJSON *AddMedicationJSONRequest

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
		UserID: auth.UserID,
		CommandBase: application.CommandBase{
			Name:                reqJSON.Name,
			InternationalName:   reqJSON.InternationalName,
			AmountValue:         reqJSON.Amount.Value,
			AmountUnit:          reqJSON.Amount.Unit,
			ReleaseForm:         reqJSON.ReleaseForm,
			Group:               reqJSON.Group,
			ManufacturerName:    reqJSON.Producer.Name,
			ManufacturerCountry: reqJSON.Producer.Country,
			ActiveSubstanceName: reqJSON.ActiveSubstance.Name,
			ActiveSubstanceDose: reqJSON.ActiveSubstance.Value,
			ActiveSubstanceUnit: reqJSON.ActiveSubstance.Unit,
			Expires:             reqJSON.Expiration,
			Release:             reqJSON.Release,
			Commentary:          reqJSON.Commentary,
		},
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
		ID: serviceResponse.ID,
		BodyCommonObject: BodyCommonObject{
			Name:              serviceResponse.Name,
			InternationalName: serviceResponse.InternationalName,
			Amount: AmountObject{
				Value: serviceResponse.AmountValue,
				Unit:  serviceResponse.AmountUnit,
			},
			ReleaseForm: serviceResponse.ReleaseForm,
			Group:       serviceResponse.Group,
			Producer: ProducerObject{
				Name:    serviceResponse.ManufacturerName,
				Country: serviceResponse.ManufacturerCountry,
			},
			ActiveSubstance: ActiveSubstanceObject{
				Name:  serviceResponse.ActiveSubstanceName,
				Value: serviceResponse.ActiveSubstanceDose,
				Unit:  serviceResponse.ActiveSubstanceUnit,
			},
			Expiration: serviceResponse.Expires,
			Release:    serviceResponse.Release,
			Commentary: serviceResponse.Commentary,
		},
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
	// embedded struct
	BodyCommonObject `json:",inline"`
}

// UpdateMedicationJSONResponse is a response for UpdateMedication.
type UpdateMedicationJSONResponse struct {
	// embedded struct
	BodyCommonObject `json:",inline"`

	ID string `json:"id"`
}

// UpdateMedication updates a medication.
func (h *MedicationHandlers) UpdateMedication(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       nil,
		})
		return
	}
	var reqJSON *UpdateMedicationJSONRequest

	vars := mux.Vars(r)
	id := vars[SlugID]

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
		UserID: auth.UserID,
		ID:     id,
		CommandBase: application.CommandBase{
			Name:                reqJSON.Name,
			InternationalName:   reqJSON.InternationalName,
			AmountValue:         reqJSON.Amount.Value,
			AmountUnit:          reqJSON.Amount.Unit,
			ReleaseForm:         reqJSON.ReleaseForm,
			Group:               reqJSON.Group,
			ManufacturerName:    reqJSON.Producer.Name,
			ManufacturerCountry: reqJSON.Producer.Country,
			ActiveSubstanceName: reqJSON.ActiveSubstance.Name,
			ActiveSubstanceDose: reqJSON.ActiveSubstance.Value,
			ActiveSubstanceUnit: reqJSON.ActiveSubstance.Unit,
			Expires:             reqJSON.Expiration,
			Release:             reqJSON.Release,
			Commentary:          reqJSON.Commentary,
		},
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
		ID: serviceResponse.ID,
		BodyCommonObject: BodyCommonObject{
			Name:              serviceResponse.Name,
			InternationalName: serviceResponse.InternationalName,
			Amount: AmountObject{
				Value: serviceResponse.AmountValue,
				Unit:  serviceResponse.AmountUnit,
			},
			ReleaseForm: serviceResponse.ReleaseForm,
			Group:       serviceResponse.Group,
			Producer: ProducerObject{
				Name:    serviceResponse.ManufacturerName,
				Country: serviceResponse.ManufacturerCountry,
			},
			ActiveSubstance: ActiveSubstanceObject{
				Name:  serviceResponse.ActiveSubstanceName,
				Value: serviceResponse.ActiveSubstanceDose,
				Unit:  serviceResponse.ActiveSubstanceUnit,
			},
			Expiration: serviceResponse.Expires,
			Release:    serviceResponse.Release,
			Commentary: serviceResponse.Commentary,
		},
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
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       nil,
		})
		return
	}
	vars := mux.Vars(r)
	id := vars[SlugID]

	serviceRequest := &application.DeleteMedicationCommand{
		UserID: auth.UserID,
		ID:     id,
	}

	_, err = h.app.DeleteMedication.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		if errors.Is(err, application.ErrDeleteInvalidUUIDFormat) {
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

// GetMedicationBoxItem returns a Box of medications.
type GetMedicationBoxItem struct {
	// embedded struct
	BodyCommonObject `json:",inline"`
}

// GetMedicationBoxJSONResponse returns a Box of medications.
type GetMedicationBoxJSONResponse struct {
	MedicationBox []GetMedicationBoxItem `json:"medicationBox"`
}

// GetMedicationBox returns a Box of medications.
func (h *MedicationHandlers) GetMedicationBox(w http.ResponseWriter, r *http.Request) {
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       nil,
		})
		return
	}
	command := &application.GetMedicationBoxCommand{
		UserID: auth.UserID,
	}
	serviceResponse, err := h.app.GetMedicationBox.Execute(
		r.Context(),
		command,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get medication Box")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedicationBox,
		})
		return
	}

	response := &GetMedicationBoxJSONResponse{
		MedicationBox: make([]GetMedicationBoxItem, 0),
	}

	for _, item := range serviceResponse.Box {
		response.MedicationBox = append(response.MedicationBox, GetMedicationBoxItem{
			BodyCommonObject: BodyCommonObject{
				Name:              item.Name,
				InternationalName: item.InternationalName,
				Amount: AmountObject{
					Value: item.AmountValue,
					Unit:  item.AmountUnit,
				},
				ReleaseForm: item.ReleaseForm,
				Group:       item.Group,
				Producer: ProducerObject{
					Name:    item.ManufacturerName,
					Country: item.ManufacturerCountry,
				},
				ActiveSubstance: ActiveSubstanceObject{
					Name:  item.ActiveSubstanceName,
					Value: item.ActiveSubstanceDose,
					Unit:  item.ActiveSubstanceUnit,
				},
				Expiration: item.Expires,
				Release:    item.Release,
				Commentary: item.Commentary,
			},
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// DataMatrixInformationJSONResponse is a response for DataMatrixInformation.
type DataMatrixInformationJSONResponse struct {
	// embedded struct
	BodyAPIObject `json:",inline"`
}

// DataMatrixInformation adds a medication.
func (h *MedicationHandlers) DataMatrixInformation(w http.ResponseWriter, r *http.Request) {
	dataParam := r.URL.Query().Get("data")
	if dataParam == "" {
		h.logger.Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgFailedToReadBody,
			Body:       nil,
		})

		return
	}

	serviceRequest := &application.DataMatrixInformationCommand{
		Data: dataParam,
	}

	serviceResponse, err := h.app.DataMatrixInformation.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		if errors.Is(err, application.ErrCantSetCache) {
			h.logger.WithError(err).Error("Failed to set cache: %w", err)
		} else {
			h.logger.WithError(err).Error("Failed to get info from scan: %w", err)
			w.WriteHeader(http.StatusInternalServerError)

			_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
				StatusCode: http.StatusInternalServerError,
				Body:       nil,
				Error:      MsgFailedToGetIfoFromScan,
			})
			return
		}
	}
	response := &DataMatrixInformationJSONResponse{
		BodyAPIObject: BodyAPIObject{
			Name:              serviceResponse.Name,
			InternationalName: serviceResponse.InternationalName,
			ReleaseForm:       serviceResponse.ReleaseForm,
			Group:             serviceResponse.Group,
			Producer: ProducerObject{
				Name:    serviceResponse.ManufacturerName,
				Country: serviceResponse.ManufacturerCountry,
			},
			Expiration: serviceResponse.Expires,
			Release:    serviceResponse.Release,
		},
	}
	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}
