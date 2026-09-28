package store

import (
	"context"
	"sync"
)

// Compile-time check: if MemoryStore stops matching Store, this line fails to compile.
var _ Store = (*MemoryStore)(nil)

// MemoryStore keeps mappings in a Go map.
//
// Interview talking points:
//   - A map is not safe for concurrent use. HTTP servers handle many
//     goroutines at once, so every read and write goes through a mutex.
//   - RWMutex: redirects are reads (RLock), creates are writes (Lock).
//     Many redirects can run together; a create waits for them to finish.
//   - This is not durable. Restart the process and every short link is gone.
//     That is why Postgres comes later. The rest of the app will not care,
//     because it only sees the Store interface.
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
