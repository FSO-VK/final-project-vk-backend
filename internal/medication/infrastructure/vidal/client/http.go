package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(config Config) *HTTPClient {
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Second
	}
	return &HTTPClient{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       config.Timeout,
		},
		config: &config,
	}
}

func (c *HTTPClient) GetInstruction(ctx context.Context, barCode string) (*vidal.ClientResponse,error) {
	url := fmt.Sprintf("%s?filter[barCode]=%s", c.config.Endpoint, barCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, vidal.ErrBadRequest
	}

	req.Header.Set("X-Token", c.config.APIToken)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, vidal.ErrBadTransport
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	var body Response
	if err := decoder.Decode(&body); err != nil {
		return nil, err
	}

	return c.handleBody(&body)
}

func (c *HTTPClient) handleBody(body *Response) (*vidal.ClientResponse, error) {
	if !body.Success {
		return nil, nil
	}

	// it is only 1 product when searching by bar code
	product := body.Products[0]
	return &vidal.ClientResponse{
		Product: product,
	}, nil
}



