package store

import (
	"context"
	"sync"
)

// Compile-time check: if MemoryStore stops matching Store, this fails to compile.
var _ Store = (*MemoryStore)(nil)

// MemoryStore keeps mappings in a Go map.
//
// Maps are not safe for concurrent use, so every read and write takes a mutex.
// Redirects use RLock so many can run together; Save uses Lock.
// Data is lost when the process exits.
type MemoryStore struct {
	mu   sync.RWMutex
	urls map[string]string // code -> long URL
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urls: make(map[string]string),
	}
}

func (s *MemoryStore) Save(_ context.Context, code, longURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.urls[code]; exists {
		return ErrConflict
	}
	s.urls[code] = longURL
	return nil
}

func (s *MemoryStore) Get(_ context.Context, code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, ok := s.urls[code]
	if !ok {
		return "", ErrNotFound
	}
	return longURL, nil
}
