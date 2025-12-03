// Package http is a package for http handlers
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/logcon"
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

	BarCode string `json:"barCode,omitempty"`
}

// AddMedicationJSONResponse is a response for AddMedication.
type AddMedicationJSONResponse struct {
	// embedded struct
	BodyCommonObject `json:",inline"`

	BarCode string `json:"barCode,omitempty"`
	ID      string `json:"id"`
}

// AddMedication adds a medication.
func (h *MedicationHandlers) AddMedication(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)
	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
		return
	}

	var reqJSON *AddMedicationJSONRequest

	var body bytes.Buffer
	_, err = body.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil {
		logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      api.MsgBadBody,
			Body:       struct{}{},
		})

		return
	}

	err = json.Unmarshal(body.Bytes(), &reqJSON)
	if err != nil {
		logger.WithError(err).Error("Failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      api.MsgBadBody,
			Body:       struct{}{},
		})
		return
	}
	logger.Debugf("request json: %+v", reqJSON)

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
			ActiveSubstance:     convertActiveSubstances(reqJSON.ActiveSubstance),
			Expires:             reqJSON.Expiration,
			Release:             reqJSON.Release,
			Commentary:          reqJSON.Commentary,
		},
		BarCode: reqJSON.BarCode,
	}

	serviceResponse, err := h.app.AddMedication.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to add medication")

		status, body := h.handleAddServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)

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
			ActiveSubstance: convertToActiveSubstanceObject(serviceResponse.ActiveSubstance),
			Expiration:      serviceResponse.Expires,
			Release:         serviceResponse.Release,
			Commentary:      serviceResponse.Commentary,
		},
		BarCode: serviceResponse.BarCode,
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
	logger := h.getLogger(r)

	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
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
		logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		})
		return
	}

	err = json.Unmarshal(body.Bytes(), &reqJSON)
	if err != nil {
		logger.WithError(err).Error("Failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
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
			ActiveSubstance:     convertActiveSubstances(reqJSON.ActiveSubstance),
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
		logger.WithError(err).Error("Failed to update medication")

		status, body := h.handleUpdateServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
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
			ActiveSubstance: convertToActiveSubstanceObject(serviceResponse.ActiveSubstance),
			Expiration:      serviceResponse.Expires,
			Release:         serviceResponse.Release,
			Commentary:      serviceResponse.Commentary,
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
	logger := h.getLogger(r)

	auth, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
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
		logger.WithError(err).Error("Failed to delete medication")

		status, body := h.handleDeleteServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
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

	BarCode string `json:"barCode,omitempty"`
	ID      string `json:"id"`
}

// GetMedicationBoxJSONResponse returns a Box of medications.
type GetMedicationBoxJSONResponse struct {
	MedicationBox []GetMedicationBoxItem `json:"medicationBox"`
}

// GetMedicationBox returns a Box of medications.
func (h *MedicationHandlers) GetMedicationBox(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)

	authorization, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
		return
	}
	command := &application.GetMedicationBoxCommand{
		UserID: authorization.UserID,
	}

	serviceResponse, err := h.app.GetMedicationBox.Execute(
		r.Context(),
		command,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to get medication Box")

		status, body := h.handlerGetServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
		return
	}

	response := &GetMedicationBoxJSONResponse{
		MedicationBox: make([]GetMedicationBoxItem, 0),
	}

	for _, medication := range serviceResponse.MedicationBox {
		response.MedicationBox = append(response.MedicationBox, GetMedicationBoxItem{
			ID: medication.ID,
			BodyCommonObject: BodyCommonObject{
				Name:              medication.Name,
				InternationalName: medication.InternationalName,
				Amount: AmountObject{
					Value: medication.AmountValue,
					Unit:  medication.AmountUnit,
				},
				ReleaseForm: medication.ReleaseForm,
				Group:       medication.Group,
				Producer: ProducerObject{
					Name:    medication.ManufacturerName,
					Country: medication.ManufacturerCountry,
				},
				ActiveSubstance: convertToActiveSubstanceObject(medication.ActiveSubstance),
				Expiration:      medication.Expires,
				Release:         medication.Release,
				Commentary:      medication.Commentary,
			},
			BarCode: medication.BarCode,
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

	BarCode string `json:"barCode,omitempty"`
}

// DataMatrixInformation adds a medication.
func (h *MedicationHandlers) DataMatrixInformation(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)

	dataParam := r.URL.Query().Get("data")
	if dataParam == "" {
		logger.Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      api.MsgBadBody,
			Body:       struct{}{},
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
		logger.WithError(err).Errorf("service: %s", err)

		status, body := h.handleDMServiceError(err)
		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
		return
	}
	response := &DataMatrixInformationJSONResponse{
		BodyAPIObject: BodyAPIObject{
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
			Expiration: serviceResponse.Expires,
			Release:    serviceResponse.Release,
		},
		BarCode: serviceResponse.BarCode,
	}
	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// GetMedicationByIDJSONResponse is a response for GetMedicationByID handler.
type GetMedicationByIDJSONResponse struct {
	BodyCommonObject `json:",inline"`

	BarCode string `json:"barCode,omitempty"`
	ID      string `json:"id"`
}

// GetMedicationByID is a handler for getting medication by its id.
func (h *MedicationHandlers) GetMedicationByID(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)

	authorization, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id := vars[SlugID]

	command := &application.GetMedicationByIDCommand{
		UserID: authorization.UserID,
		ID:     id,
	}

	medication, err := h.app.GetMedicationByID.Execute(r.Context(), command)
	if err != nil {
		logger.WithError(err).Error("Failed to get medication by id")

		status, body := h.handleGetByIDServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
		return
	}

	response := &GetMedicationByIDJSONResponse{
		ID: medication.ID,
		BodyCommonObject: BodyCommonObject{
			Name:              medication.Name,
			InternationalName: medication.InternationalName,
			Amount: AmountObject{
				Value: medication.AmountValue,
				Unit:  medication.AmountUnit,
			},
			ReleaseForm: medication.ReleaseForm,
			Group:       medication.Group,
			Producer: ProducerObject{
				Name:    medication.ManufacturerName,
				Country: medication.ManufacturerCountry,
			},
			ActiveSubstance: convertToActiveSubstanceObject(medication.ActiveSubstance),
			Expiration:      medication.Expires,
			Release:         medication.Release,
			Commentary:      medication.Commentary,
		},
		BarCode: medication.BarCode,
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// InstructionAssistantJSONResponse is a response for InstructionAssistant.
type InstructionAssistantJSONResponse struct {
	LLMAnswer string `json:"llmAnswer"`
}

// InstructionAssistant gives advices about medication instructions.
func (h *MedicationHandlers) InstructionAssistant(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)

	authorization, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id := vars[SlugID]

	question := r.URL.Query().Get("question")
	if question == "" {
		logger.Error("Failed to read request query parameters")
		w.WriteHeader(http.StatusBadRequest)

		_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      api.MsgBadBody,
			Body:       struct{}{},
		})

		return
	}

	serviceRequest := &application.InstructionAssistantCommand{
		UserQuestion: question,
		MedicationID: id,
		UserID:       authorization.UserID,
	}

	serviceResponse, err := h.app.InstructionAssistant.Execute(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		logger.WithError(err).Errorf("service: %s", err)

		status, body := h.handleAssistantServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
		return
	}
	response := &InstructionAssistantJSONResponse{
		LLMAnswer: serviceResponse.LLMAnswer,
	}
	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// GetInstructionJSONResponse is a response for GetInstructionByMedicationID handler.
type GetInstructionJSONResponse struct {
	InstructionCommonObject `json:",inline"`
}

// GetInstruction is a handler for getting medication's instruction.
func (h *MedicationHandlers) GetInstruction(w http.ResponseWriter, r *http.Request) {
	logger := h.getLogger(r)

	authorization, err := httputil.GetAuthFromCtx(r)
	if err != nil {
		h.writeResponseUnauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id := vars[SlugID]

	command := &application.GetInstructionByMedicationIDCommand{
		UserID: authorization.UserID,
		ID:     id,
	}

	medication, err := h.app.GetInstructionByMedicationID.Execute(r.Context(), command)
	if err != nil {
		logger.WithError(err).Error("Failed to get instruction by medication id")

		status, body := h.handleGetInstructionByIDServiceError(err)

		w.WriteHeader(status)
		_ = httputil.NetHTTPWriteJSON(w, body)
		return
	}

	response := &GetInstructionJSONResponse{
		InstructionCommonObject: InstructionCommonObject{
			Nosologies:             convertToNosology(medication.Nosologies),
			ClPhPointers:           convertToClPhPointers(medication.ClPhPointers),
			PharmInfluence:         medication.PharmInfluence,
			PharmKinetics:          medication.PharmKinetics,
			Dosage:                 medication.Dosage,
			OverDosage:             medication.OverDosage,
			Interaction:            medication.Interaction,
			Lactation:              medication.Lactation,
			SideEffects:            medication.SideEffects,
			UsingIndication:        medication.UsingIndication,
			UsingCounterIndication: medication.UsingCounterIndication,
			SpecialInstruction:     medication.SpecialInstruction,
			RenalInsuf:             medication.RenalInsuf,
			HepatoInsuf:            medication.HepatoInsuf,
			ElderlyInsuf:           medication.ElderlyInsuf,
			ChildInsuf:             medication.ChildInsuf,
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// handleAddServiceError maps service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handleAddServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToAddMedication,
		}
	}
}

// handleUpdateServiceError maps service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handleUpdateServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrNoMedication):
		return http.StatusNotFound, &api.Response[any]{
			StatusCode: http.StatusNotFound,
			Body:       struct{}{},
			Error:      MsgNoMedication,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToUpdateMedication,
		}
	}
}

// handleDeleteServiceError maps service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handleDeleteServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrNoMedication):
		return http.StatusNotFound, &api.Response[any]{
			StatusCode: http.StatusNotFound,
			Body:       struct{}{},
			Error:      MsgNoMedication,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToDeleteMedication,
		}
	}
}

// handlerGetServiceError maps service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handlerGetServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedicationBox,
		}
	}
}

func (h *MedicationHandlers) handleGetByIDServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrNoMedication):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgNoMedication,
		}
	case errors.Is(err, application.ErrFailedToGetMedication):
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedication,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      api.MsgServerError,
		}
	}
}

// handleDMServiceError maps DataMatrix service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handleDMServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrCantSetCache):
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      api.MsgServerError,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetIfoFromScan,
		}
	}
}

// handleAssistantServiceError maps service errors to HTTP status and API responses using switch.
func (h *MedicationHandlers) handleAssistantServiceError(err error) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrNoMedication):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedication,
		}
	case errors.Is(err, application.ErrNoInstruction):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgFailedToGetInstructions,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetInfoFromLLM,
		}
	}
}

func (h *MedicationHandlers) handleGetInstructionByIDServiceError(
	err error,
) (int, *api.Response[any]) {
	switch {
	case errors.Is(err, application.ErrValidationFail):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		}
	case errors.Is(err, application.ErrNoMedication):
		return http.StatusBadRequest, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      MsgNoMedication,
		}
	case errors.Is(err, application.ErrFailedToGetMedication):
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetMedication,
		}
	case errors.Is(err, application.ErrFailedToGetInstruction):
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetInstructions,
		}
	default:
		return http.StatusInternalServerError, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      api.MsgServerError,
		}
	}
}

// getLogger returns a logger from the context if exists,
// otherwise returns a default logger.
func (h *MedicationHandlers) getLogger(r *http.Request) *logrus.Entry {
	l, ok := logcon.FromContext(r.Context())
	if !ok {
		return h.logger
	}
	return l
}

// writeResponseUnauthorized writes an Unauthorized HTTP response.
func (h *MedicationHandlers) writeResponseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_ = httputil.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusUnauthorized,
		Body:       struct{}{},
		Error:      api.MsgUnauthorized,
	})
}

func convertToActiveSubstanceObject(
	substances []application.ActiveSubstance,
) []ActiveSubstanceObject {
	result := make([]ActiveSubstanceObject, len(substances))
	for i, v := range substances {
		result[i] = ActiveSubstanceObject{
			Name:  v.Name,
			Value: v.Value,
			Unit:  v.Unit,
		}
	}
	return result
}

func convertActiveSubstances(substances []ActiveSubstanceObject) []application.ActiveSubstance {
	result := make([]application.ActiveSubstance, len(substances))
	for i, v := range substances {
		result[i] = application.ActiveSubstance(v)
	}
	return result
}

func convertToNosology(substances []application.Nosology) []Nosology {
	result := make([]Nosology, len(substances))
	for i, v := range substances {
		result[i] = Nosology{
			Code: v.Code,
			Name: v.Name,
		}
	}
	return result
}

func convertToClPhPointers(substances []application.ClPhPointer) []ClPhPointer {
	result := make([]ClPhPointer, len(substances))
	for i, v := range substances {
		result[i] = ClPhPointer{
			Code: v.Code,
			Name: v.Name,
		}
	}
	return result
}
