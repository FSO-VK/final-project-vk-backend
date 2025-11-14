package vidal

import (
	"context"
	"errors"
)

type Client interface {
	GetInstruction(ctx context.Context, barCode string) (*ClientResponse, error)
}

var (
	ErrBadRequest = errors.New("bad request")
	ErrBadTransport  = errors.New("broken transport")
)

type ClientResponse struct {
	Product
}
