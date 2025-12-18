package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	strip "github.com/grokify/html-strip-tags-go"
)

// HTTPClient is a client for vidal.ru API.
type HTTPClient struct {
	client *http.Client
	config *Config
}

// NewHTTPClient creates a new HTTPClient.
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

// GetInstruction returns a product info by bar code.
func (c *HTTPClient) GetInstruction(
	ctx context.Context,
	barCode string,
) (*vidal.ClientResponse, error) {
	url := fmt.Sprintf("%s?filter[barCode]=%s", c.config.Endpoint, barCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Token", c.config.APIToken)
	log.Printf("response %+v", req)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform HTTP request: %w", err)
	}
	log.Printf("response %+v", res)

	if res.StatusCode != http.StatusOK {
		//nolint:err113 // need to return status code, so dynamic errors is more readable.
		return nil, fmt.Errorf("got unexpected status code: %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()

	var body Response
	if err := decoder.Decode(&body); err != nil {
		return nil, fmt.Errorf("decode JSON response: %w", err)
	}
	log.Printf("response body %+v", body)

	return c.handleBody(&body)
}

func (c *HTTPClient) handleBody(body *Response) (*vidal.ClientResponse, error) {
	if !body.Success {
		return nil, fmt.Errorf(
			"json response status is unsuccessful: %w",
			vidal.ErrStorageNoProduct,
		)
	}

	if len(body.Products) == 0 {
		return nil, vidal.ErrClientNoProduct
	}

	// it is only 1 possible product when searching by bar code
	product := body.Products[0]

	// vidal API contains HTML tags in document (instruction) fields.
	err := stripHTMLTags(&product.Document)
	if err != nil {
		return nil, fmt.Errorf("strip HTML tags: %w", err)
	}

	return &vidal.ClientResponse{
		Product: product,
	}, nil
}

var errNilDocument = errors.New("document is nil")

func stripHTMLTags(doc *vidal.Document) error {
	if doc == nil {
		return errNilDocument
	}

	doc.PhInfluence = strip.StripTags(doc.PhInfluence)
	doc.PhKinetics = strip.StripTags(doc.PhKinetics)
	doc.Dosage = strip.StripTags(doc.Dosage)
	doc.OverDosage = strip.StripTags(doc.OverDosage)
	doc.Interaction = strip.StripTags(doc.Interaction)
	doc.Lactation = strip.StripTags(doc.Lactation)
	doc.SideEffects = strip.StripTags(doc.SideEffects)
	doc.Indication = strip.StripTags(doc.Indication)
	doc.ContraIndication = strip.StripTags(doc.ContraIndication)
	doc.SpecialInstruction = strip.StripTags(doc.SpecialInstruction)
	doc.PregnancyUsing = strip.StripTags(doc.PregnancyUsing)
	doc.NursingUsing = strip.StripTags(doc.NursingUsing)
	doc.RenalInsuf = strip.StripTags(doc.RenalInsuf)
	doc.HepatoInsuf = strip.StripTags(doc.HepatoInsuf)
	doc.PharmDelivery = strip.StripTags(doc.PharmDelivery)
	doc.ElderlyInsuf = strip.StripTags(doc.ElderlyInsuf)
	doc.ChildInsuf = strip.StripTags(doc.ChildInsuf)

	return nil
}
