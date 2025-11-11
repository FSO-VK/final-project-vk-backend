package llmclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

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
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	credentials := clientID + ":" + clientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(credentials))

	req.Header.Set("Authorization", "Basic "+basicAuth)
	//nolint:canonicalheader
	req.Header.Set("RqUID", uuid.New().String()) // ну а у гигачата вот такой хедер неправильный
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := createHTTPClient(timeout)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
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
