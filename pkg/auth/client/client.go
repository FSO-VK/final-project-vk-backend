// Package auth provides auth client for other services.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// HTTPAuthChecker implements AuthChecker interface.
type HTTPAuthChecker struct {
	client *http.Client
	cfg    ClientConfig
	logger *logrus.Entry
}

// NewHTTPAuthChecker creates a new HTTPAuthChecker.
func NewHTTPAuthChecker(cfg ClientConfig, logger *logrus.Entry) *HTTPAuthChecker {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	return &HTTPAuthChecker{client: client, cfg: cfg, logger: logger}
}

// CheckAuth checks auth for other services.
func (h *HTTPAuthChecker) CheckAuth(reqData *Request) (*Response, error) {
	if err := h.checkRequest(reqData); err != nil {
		return nil, err
	}

	ctx := context.Background()

	if h.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.cfg.Timeout)
		defer cancel()
	}

	url := h.cfg.BaseURL + h.cfg.Path
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.AddCookie(&http.Cookie{
		Name:  h.cfg.CookieName,
		Value: reqData.SessionID,
		Path:  h.cfg.CookieDomain,
	})
	httpReq.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		h.logger.WithError(err).Warn("auth request failed")
		return nil, ErrAuthServiceUnavailable
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.WithError(err).Debug("failed to close response body")
		}
	}()

	if body, err := h.checkHTTPStatus(resp.StatusCode); err != nil {
		return body, err
	}

	var env api.Response[ExpectedCheckAuthResponse]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		h.logger.WithError(err).Error("failed to decode auth response")
		return nil, ErrInvalidAuthResponse
	}

	if env.StatusCode == http.StatusUnauthorized || env.StatusCode == http.StatusForbidden {
		return &Response{
			SessionID:    reqData.SessionID,
			UserID:       "",
			IsAuthorized: false,
		}, nil
	}

	out := &Response{
		SessionID:    reqData.SessionID,
		UserID:       env.Body.UserID,
		IsAuthorized: true,
	}

	return out, nil
}

func (h *HTTPAuthChecker) checkRequest(reqData *Request) error {
	if reqData == nil || reqData.SessionID == "" {
		return ErrInvalidRequest
	}
	return nil
}

func (h *HTTPAuthChecker) checkHTTPStatus(statusCode int) (*Response, error) {
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return &Response{
			SessionID:    "",
			UserID:       "",
			IsAuthorized: false,
		}, nil
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, ErrBadResponse
	}
	return &Response{}, nil
}
