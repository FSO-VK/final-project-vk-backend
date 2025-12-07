// Package medicationclient implements MedicationService interface for getting medication info from medication service.
package medicationclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MedicationClient implements MedicationService.
type MedicationClient struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

// NewMedicationClient creates a new MedicationClient.
func NewMedicationClient(cfg ClientConfig, logger *logrus.Entry) *MedicationClient {
	client := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
	return &MedicationClient{client: client, cfg: cfg, logger: logger}
}

type MedicationExpectedResponse struct {
	StatusCode int  `json:"statusCode"`
	Body       Body `json:"body"`
}

type Body struct {
	Name string `json:"name"`
}

// MedicationName implements MedicationClient interface.
func (h *MedicationClient) MedicationName(
	id uuid.UUID,
) (string, error) {
	parsedResponse, err := h.makeFullRequest(id)
	if err != nil {
		return "", err
	}

	return parsedResponse.Name, nil
}

func (h *MedicationClient) makeFullRequest(id uuid.UUID) (Body, error) {
	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}

	url := h.cfg.Endpoint + id.String()
	httpReq, err := http.NewRequestWithContext(ctx, h.cfg.Method, url, nil)
	if err != nil {
		return Body{}, err
	}

	resp, err := h.client.Do(httpReq)
	if err != nil {
		h.logger.WithError(err).Warn("medication API request failed")
		return Body{}, ErrMedicationServiceUnavailable
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return Body{}, ErrBadResponse
	}

	var parsedResponse MedicationExpectedResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResponse); err != nil {
		h.logger.WithError(err).Error("failed to decode medication API response")
		return Body{}, fmt.Errorf("%w: %w", ErrBadResponse, err)
	}

	if parsedResponse.StatusCode != http.StatusOK {
		return Body{}, ErrNoMedicationFound
	}

	return parsedResponse.Body, nil
}
