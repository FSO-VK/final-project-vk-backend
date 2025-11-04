// Package datamatrix implements DataMatrixClient interface for getting medication info from dataMatrix API.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	clientInterface "github.com/FSO-VK/final-project-vk-backend/internal/medication/application/api_client"
	"github.com/sirupsen/logrus"
)

// APIDataMatrix implements AuthChecker DataMatrixClient.
type APIDataMatrix struct {
	client *http.Client
	cfg    PushClient
	logger *logrus.Entry
}

// NewDataMatrixAPI creates a new DataMatrixAPI.
func NewDataMatrixAPI(cfg PushClient, logger *logrus.Entry) *APIDataMatrix {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &APIDataMatrix{client: client, cfg: cfg, logger: logger}
}

// GetInformationByDataMatrix implements DataMatrixClient interface.
func (h *APIDataMatrix) GetInformationByDataMatrix(
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, ErrBadResponse
	}

	var parsedResponse ExpectedDataMatrixAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResponse); err != nil {
		h.logger.WithError(err).Error("failed to decode API response")
		return nil, fmt.Errorf("%w: %w", ErrBadResponse, err)
	}

	if !parsedResponse.CodeFound {
		return nil, ErrNoMedicationFound
	}

	out := MapToMedicationInfo(&parsedResponse)

	return out, nil
}

func (h *APIDataMatrix) checkRequest(data *clientInterface.DataMatrixCodeInfo) error {
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

	firstRune, size := utf8.DecodeRuneInString(name)

	return string(unicode.ToUpper(firstRune)) + strings.ToLower(name[size:])
}
