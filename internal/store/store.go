package store

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a short code has no mapping.
// Handlers can map this to HTTP 404 without knowing about the database.
var ErrNotFound = errors.New("short url not found")

// Store is the persistence port for short-code -> long-URL mappings.
//
// Interview talking point: HTTP and business logic depend on this interface,
// not on Postgres, Redis, or a map. That is dependency inversion.
// We can start with memory, then swap in Postgres, without rewriting handlers.
//
// Save and Get take context.Context so timeouts and cancellation from the
// HTTP request flow through to storage. That matters once we leave in-memory.
type Store interface {
	Save(ctx context.Context, code, longURL string) error
	Get(ctx context.Context, code string) (string, error)
}
