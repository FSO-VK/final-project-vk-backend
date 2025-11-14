package vidal

import (
	"context"
	"errors"
)

// ErrClientNoProduct means that a product is not found in external service.
var ErrClientNoProduct = errors.New("no product found")

// Client is an interface for client to external service.
type Client interface {
	GetInstruction(ctx context.Context, barCode string) (*ClientResponse, error)
}

// ClientResponse is a response from external service.
type ClientResponse struct {
	Product
}
