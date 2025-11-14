package vidal

import "context"

type Storage interface {
	SaveProductInfo(ctx context.Context, product *Product) error
	GetProductInfo(ctx context.Context, barCode string) (*Product, error)
}

type StorageModel struct {
	Product

	BarCodes []string
}
