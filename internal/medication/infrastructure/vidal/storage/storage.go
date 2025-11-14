// Package storage contains cache for in memory db.
// It's only for PR, after review it will be changed to a MongoDB.
package storage

import (
	"context"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
)

// ProductCache implements a storage for API's responses.
type ProductCache struct {
	data *cache.Cache[*vidal.StorageModel]
	mu   *sync.RWMutex
}

// NewProductCache creates a ProductCache.
func NewProductCache() *ProductCache {
	return &ProductCache{
		data: cache.NewCache[*vidal.StorageModel](),
		mu:   &sync.RWMutex{},
	}
}

// SaveProduct sets a product in storage.
func (p *ProductCache) SaveProduct(_ context.Context, product *vidal.StorageModel) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, barcode := range product.BarCodes {
		p.data.Set(barcode, product)
	}
	return nil
}

// GetProduct returns a product from storage.
func (p *ProductCache) GetProduct(_ context.Context, barCode string) (*vidal.StorageModel, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if product, ok := p.data.Get(barCode); ok {
		return product, nil
	}
	return nil, vidal.ErrStorageNoProduct
}
