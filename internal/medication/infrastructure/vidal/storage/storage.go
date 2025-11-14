package storage

import (
	"context"
	"sync"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
)

type ProductCache struct {
	data *cache.Cache[*vidal.StorageModel]
	mu   *sync.RWMutex
}

func NewProductCache() *ProductCache {
	return &ProductCache{
		data: cache.NewCache[*vidal.StorageModel](),
		mu:   &sync.RWMutex{},
	}
}

func (p *ProductCache) SetProduct(_ context.Context, product *vidal.StorageModel) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, barcode := range product.BarCodes {
		p.data.Set(barcode, product)
	}
	return nil
}

func (p *ProductCache) GetProduct(_ context.Context, barCode string) (*vidal.StorageModel, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if product, ok := p.data.Get(barCode); ok {
		return product, nil
	}
	return nil, vidal.ErrNoProduct
}
