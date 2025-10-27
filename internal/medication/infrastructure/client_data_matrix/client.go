// Package clientdatamatrix implements DataMatrixClient interface for getting medication info from dataMatrix API.
package clientdatamatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"unicode"

	clientInterface "github.com/FSO-VK/final-project-vk-backend/internal/medication/application/api_client"
	"github.com/sirupsen/logrus"
)

// DataMatrixAPI implements AuthChecker DataMatrixClient.
type DataMatrixAPI struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

// NewDataMatrixAPI creates a new DataMatrixAPI.
func NewDataMatrixAPI(cfg ClientConfig, logger *logrus.Entry) *DataMatrixAPI {
	client := &http.Client{
		Timeout:       2 * time.Second,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &DataMatrixAPI{client: client, cfg: cfg, logger: logger}
}

// GetInformationByDataMatrix implements DataMatrixClient interface.
func (h *DataMatrixAPI) GetInformationByDataMatrix(
	data *clientInterface.DataMatrixCodeInfo,
) (*clientInterface.MedicationInfo, error) {
	if err := h.checkRequest(data); err != nil {
		return nil, err
	}

	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}
	code := "01" + data.GTIN + "21" + data.SerialNumber + "%1D91" + data.CryptoData91 + "%1D92" + data.CryptoData92
	url := h.cfg.Host + code
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(httpReq)
	if err != nil {
		h.logger.WithError(err).Warn("dataMatrix API request failed")
		return nil, ErrAuthServiceUnavailable
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.WithError(err).Debug("failed to close response body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ErrInvalidAPIResponse
	}

	var env ExpectedDataMatrixAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		h.logger.WithError(err).Error("failed to decode API response")
		return nil, ErrInvalidAPIResponse
	}

	if !env.CodeFounded {
		return nil, ErrNoMedicationFound
	}

	out := MapToMedicationInfo(&env)

	return out, nil
}

func (h *DataMatrixAPI) checkRequest(data *clientInterface.DataMatrixCodeInfo) error {
	if data == nil || data.GTIN == "" || data.SerialNumber == "" ||
		data.CryptoData91 == "" || data.CryptoData92 == "" {
		return ErrInvalidRequest
	}
	return nil
}

// MapToMedicationInfo maps ExpectedDataMatrixAPIResponse to MedicationInfo.
func MapToMedicationInfo(resp *ExpectedDataMatrixAPIResponse) *clientInterface.MedicationInfo {
	if resp == nil {
		return nil
	}

	return &clientInterface.MedicationInfo{
		Name:                resp.ProductName,
		InternationalName:   resp.DrugsData.FOIV.ProductName,
		AmountValue:         0,
		AmountUnit:          "",
		ReleaseForm:         formatReleaseForm(resp.DrugsData.FOIV.ProductFormName),
		Group:               resp.Category,
		ManufacturerName:    resp.DrugsData.FOIV.Manufacturer,
		ManufacturerCountry: resp.DrugsData.FOIV.ManufacturerCountry,
		ActiveSubstanceName: "",
		ActiveSubstanceDose: 0,
		ActiveSubstanceUnit: "",
		Expires:             resp.ExpDate,
		Release:             formatReleaseDate(resp.DrugsData.ReleaseDate),
	}
}

func formatReleaseDate(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp/1000, 0).Format("2006-01-02")
}

func formatReleaseForm(name string) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	if len(runes) == 0 {
		return ""
	}

	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}

	return string(runes)
}
