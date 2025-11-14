package vidal

import (
	"context"
	"errors"
)

var ErrNoProduct = errors.New("no product found")

type Storage interface {
	SaveProduct(ctx context.Context, product *StorageModel) error
	GetProduct(ctx context.Context, barCode string) (*StorageModel, error)
}

type StorageModel struct {
	Product

	BarCodes []string
}
