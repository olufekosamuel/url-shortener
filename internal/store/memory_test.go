package store

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestMemoryStoreSaveAndGet(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	if err := s.Save(ctx, "abc123", "https://example.com"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Get(ctx, "abc123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "https://example.com" {
		t.Fatalf("Get = %q, want https://example.com", got)
	}
}

func TestMemoryStoreMissingCode(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.Get(context.Background(), "nope")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get missing: err = %v, want ErrNotFound", err)
	}
}

func TestMemoryStoreConflict(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	if err := s.Save(ctx, "abc123", "https://example.com/a"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	err := s.Save(ctx, "abc123", "https://example.com/b")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("second Save: err = %v, want ErrConflict", err)
	}
}

// Without the mutex, concurrent map access can panic.
func TestMemoryStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		code := fmt.Sprintf("c%d", i)
		url := fmt.Sprintf("https://example.com/%d", i)
		wg.Add(2)

		go func(code, url string) {
			defer wg.Done()
			if err := s.Save(ctx, code, url); err != nil {
				t.Errorf("Save: %v", err)
			}
		}(code, url)

		go func(code string) {
			defer wg.Done()
			_, _ = s.Get(ctx, code)
		}(code)
	}
	wg.Wait()
}
