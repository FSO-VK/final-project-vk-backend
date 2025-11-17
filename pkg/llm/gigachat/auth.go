package gigachat

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrAuthFailed is returned when the LLM server returns an unsuccessful auth status.
	ErrAuthFailed = errors.New("auth failed")
	// ErrEmptyAccessToken is returned when the access token is missing in the response.
	ErrEmptyAccessToken = errors.New("empty access token in response")
	// ErrBadRequestData is returned when the request data is invalid.
	ErrBadRequestData = errors.New("bad request data")
)

// GigaChatAuth is the response body for the GigaChat API.
type GigaChatAuth struct {
	//nolint:tagliatelle
	AccessToken string `json:"access_token"`
	//nolint:tagliatelle
	ExpiresAt int64 `json:"expires_at"`
}

func createHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
}

func getGigaChatToken(
	clientID string,
	clientSecret string,
	authURL string,
	timeout time.Duration,
) (string, error) {
	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		authURL,
		bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", ErrBadRequestData
	}

	credentials := clientID + ":" + clientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(credentials))

	req.Header.Set("Authorization", "Basic "+basicAuth)
	//nolint:canonicalheader // RqUID is sber header we need to send
	req.Header.Set("RqUID", uuid.New().String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := createHTTPClient(timeout)

	resp, err := client.Do(req)
	if err != nil {
		return "", ErrAuthFailed
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrAuthFailed
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrAuthFailed
	}

	var authResp GigaChatAuth
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", ErrAuthFailed
	}

	if authResp.AccessToken == "" {
		return "", ErrAuthFailed
	}

	return authResp.AccessToken, nil
}
