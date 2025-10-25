package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/sirupsen/logrus"
)

type HTTPAuthChecker struct {
	doer   HttpDoer
	cfg    ClientConfig
	logger *logrus.Entry
}

func NewHTTPAuthChecker(doer HttpDoer, cfg ClientConfig, logger *logrus.Entry) *HTTPAuthChecker {
	return &HTTPAuthChecker{doer: doer, cfg: cfg, logger: logger}
}

func (h *HTTPAuthChecker) CheckAuth(reqData *Request) (*Response, error) {
	if reqData == nil || reqData.SessionID == "" {
		return &Response{
			SessionID:    "",
			UserID:       "",
			IsAuthorized: false,
		}, nil
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

	resp, err := h.doer.Do(httpReq)
	if err != nil {
		h.logger.WithError(err).Warn("auth request failed")
		return nil, ErrAuthServiceUnavailable
	}
	defer h.safeClose(resp.Body)

	if err := h.checkHTTPStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	var env api.Response[ExpectedCheckAuthResponse]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		h.logger.WithError(err).Error("failed to decode auth response")
		return nil, ErrUnauthorized
	}

	if env.StatusCode == http.StatusUnauthorized || env.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	out := &Response{
		SessionID:    reqData.SessionID,
		UserID:       env.Body.UserID,
		IsAuthorized: true,
	}

	return out, nil
}

func (h *HTTPAuthChecker) safeClose(body io.Closer) {
	_ = body.Close()
}

func (h *HTTPAuthChecker) checkHTTPStatus(statusCode int) error {
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	if statusCode < 200 || statusCode >= 300 {
		return ErrBadResponse
	}
	return nil
}
