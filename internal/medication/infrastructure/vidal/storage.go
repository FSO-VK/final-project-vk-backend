package vidal

import (
	"context"
	"errors"
	"time"
)

// ErrStorageNoProduct occurs when a product is not found in storage.
var ErrStorageNoProduct = errors.New("no product found")

// Storage is an interface for storage.
type Storage interface {
	SaveProduct(ctx context.Context, product *StorageModel) error
	GetProduct(ctx context.Context, barCode string) (*StorageModel, error)
}

// StorageModel is a model for storage.
type StorageModel struct {
	Product

	BarCodes  []string  `bson:"barCodes"`
	CreatedAt time.Time `bson:"createdAt"`
}
