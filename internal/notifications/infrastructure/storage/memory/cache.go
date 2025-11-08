// Package memory contains implementation of medication's repositories.
package memory

import (
	"sync"
)

// Cache is a cache realization for in memory db.
type Cache[T any] struct {
	data map[string]T
	size uint64

	mu *sync.RWMutex
}

// NewCache returns a new cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		data: make(map[string]T),
		size: 0,
		mu:   &sync.RWMutex{},
	}
}

// Set sets a value in the cache.
func (s *Cache[T]) Set(key string, value T) {
	s.mu.Lock()
	s.data[key] = value
	s.size++
	s.mu.Unlock()
}

// Get returns a value from the cache.
func (s *Cache[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return s.data[key], ok
}

// Delete deletes a value from the cache.
func (s *Cache[T]) Delete(key string) {
	s.mu.Lock()
	delete(s.data, key)
	s.size--
	s.mu.Unlock()
}
