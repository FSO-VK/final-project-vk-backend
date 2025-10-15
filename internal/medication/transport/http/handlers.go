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
	SlugID = "id"
)

type MedicineHandlers struct {
	MedicineService application.MedicineService
	logger          *logrus.Entry
}

func NewHandlers(
	medicineService application.MedicineService,
	logger *logrus.Entry,
) *MedicineHandlers {
	return &MedicineHandlers{
		MedicineService: medicineService,
		logger:          logger,
	}
}

type AddMedicineJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

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
			Error:      "Failed to read request body",
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
			Error:      "Failed to unmarshal request body",
			Body:       nil,
		})
		return
	}

	serviceRequest := &application.AddMedicineRequest{
		Name:      reqJSON.Name,
		Items:     reqJSON.Items,
		ItemsUnit: reqJSON.ItemsUnit,
		Expires:   reqJSON.Expiration,
	}

	serviceResponse, err := h.MedicineService.AddMedicine(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       nil,
			Error:      "Failed to add medicine",
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       serviceResponse,
	})
}

type UpdateMedicineJSONRequest struct {
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

type UpdateMedicineJSONResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

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
			Error:      "Failed to parse id",
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
			Error:      "Failed to read request body",
		})
		return
	}

	err = json.Unmarshal(body.Bytes(), &reqJSON)
	if err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      "Failed to unmarshal request body",
		})
		return
	}

	serviceRequest := &application.UpdateMedicineRequest{
		ID:        uint(idUint),
		Name:      reqJSON.Name,
		Items:     reqJSON.Items,
		ItemsUnit: reqJSON.ItemsUnit,
		Expires:   reqJSON.Expiration,
	}

	serviceResponse, err := h.MedicineService.UpdateMedicine(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed to update medicine",
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
	})
}

func (h *MedicineHandlers) DeleteMedicine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[SlugID]
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to parse id")
		w.WriteHeader(http.StatusBadRequest)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      "Failed to parse id",
		})
		return
	}

	serviceRequest := &application.DeleteMedicineRequest{
		ID: uint(idUint),
	}

	_, err = h.MedicineService.DeleteMedicine(
		r.Context(),
		serviceRequest,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete medicine")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed to delete medicine",
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = httph.NetHTTPWriteJSON(w, &api.Response[struct{}]{
		StatusCode: http.StatusOK,
		Body:       struct{}{},
	})
}

type GetMedicineListItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Items      uint   `json:"items"`
	ItemsUnit  string `json:"itemsUnit"`
	Expiration string `json:"expiration"`
}

type GetMedicineListJSONResponse struct {
	MedicineList []GetMedicineListItem `json:"medicineList"`
}

func (h *MedicineHandlers) GetMedicineList(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := h.MedicineService.GetMedicineList(
		r.Context(),
		&application.GetMedicineListRequest{},
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get medicine list")
		w.WriteHeader(http.StatusInternalServerError)

		_ = httph.NetHTTPWriteJSON(w, &api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed to get medicine list",
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
	})
}
