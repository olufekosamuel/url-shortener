package store

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is returned when a short code has no mapping.
	ErrNotFound = errors.New("short url not found")
	// ErrConflict is returned when Save is asked to use a code that already exists.
	ErrConflict = errors.New("short code already exists")
)

// Store persists short-code → long-URL mappings.
//
// HTTP and domain logic depend on this interface, not on a map or a database.
// Save and Get take context.Context so request timeouts can reach storage
// once a remote backend is used.
type Store interface {
	Save(ctx context.Context, code, longURL string) error
	Get(ctx context.Context, code string) (string, error)
}
