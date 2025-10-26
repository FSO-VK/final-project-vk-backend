package datamatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// DataMatrixAPI implements AuthChecker DataMatrixClient.
type DataMatrixAPI struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

func NewDataMatrixAPI(cfg ClientConfig, logger *logrus.Entry) *DataMatrixAPI {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	return &DataMatrixAPI{client: client, cfg: cfg, logger: logger}
}

func (h *DataMatrixAPI) GetInformationByDataMatrix(data *DataMatrixScannedInfo) (*MedicationInfoFromAPI, error) {
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
	url := h.cfg.BaseURL + code

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

	if env.CodeFounded == false {
		return nil, ErrNoMedicationFound
	}

	out := &MedicationInfoFromAPI{
		ExpDate: env.ExpDate,
	}

	return out, nil
}

func (h *DataMatrixAPI) checkRequest(data *DataMatrixScannedInfo) error {
	if data == nil || data.GTIN == "" || data.SerialNumber == "" || data.CryptoData91 == "" || data.CryptoData92 == "" {
		return ErrInvalidRequest
	}
	return nil
}
