package memory

import (
	"sync"
)

type Cache[T any] struct {
	data map[string]T
	size uint64

	mu *sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		data: make(map[string]T),
		size: 0,
		mu:   &sync.RWMutex{},
	}
}

func (s *Cache[T]) Set(key string, value T) {
	s.mu.Lock()
	s.data[key] = value
	s.size++
	s.mu.Unlock()
}

func (s *Cache[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return s.data[key], ok
}

func (s *Cache[T]) Delete(key string) {
	s.mu.Lock()
	delete(s.data, key)
	s.size--
	s.mu.Unlock()
}
