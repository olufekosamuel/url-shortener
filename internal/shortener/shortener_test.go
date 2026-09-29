package shortener

import (
	"context"
	"errors"
	"testing"
	"unicode"

	"github.com/olufekosamuel/url-shortener/internal/store"
)

func TestShortenAndResolve(t *testing.T) {
	svc := New(store.NewMemoryStore())
	svc.generate = func() (string, error) { return "abc1234", nil }

	code, err := svc.Shorten(context.Background(), "https://example.com/a")
	if err != nil {
		t.Fatalf("Shorten: %v", err)
	}
	if code != "abc1234" {
		t.Fatalf("Shorten code = %q, want abc1234", code)
	}

	got, err := svc.Resolve(context.Background(), code)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got != "https://example.com/a" {
		t.Fatalf("Resolve = %q, want https://example.com/a", got)
	}
}

func TestShortenRejectsInvalidURL(t *testing.T) {
	svc := New(store.NewMemoryStore())
	ctx := context.Background()

	for _, raw := range []string{"", "example.com", "javascript:alert(1)", "ftp://example.com"} {
		_, err := svc.Shorten(ctx, raw)
		if !errors.Is(err, ErrInvalidURL) {
			t.Errorf("Shorten(%q): err = %v, want ErrInvalidURL", raw, err)
		}
	}
}

func TestShortenRetriesOnConflict(t *testing.T) {
	st := store.NewMemoryStore()
	ctx := context.Background()
	if err := st.Save(ctx, "aaaaaaa", "https://example.com/old"); err != nil {
		t.Fatalf("seed Save: %v", err)
	}

	n := 0
	svc := New(st)
	svc.generate = func() (string, error) {
		n++
		if n == 1 {
			return "aaaaaaa", nil
		}
		return "bbbbbbb", nil
	}

	code, err := svc.Shorten(ctx, "https://example.com/new")
	if err != nil {
		t.Fatalf("Shorten: %v", err)
	}
	if code != "bbbbbbb" {
		t.Fatalf("Shorten code = %q, want bbbbbbb after retry", code)
	}
}

func TestResolveMissing(t *testing.T) {
	svc := New(store.NewMemoryStore())
	_, err := svc.Resolve(context.Background(), "missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("Resolve: err = %v, want ErrNotFound", err)
	}
}

func TestRandomCode(t *testing.T) {
	got, err := RandomCode(7)
	if err != nil {
		t.Fatalf("RandomCode: %v", err)
	}
	if len(got) != 7 {
		t.Fatalf("len = %d, want 7", len(got))
	}
	for _, r := range got {
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			t.Fatalf("code %q contains %q, want Base62", got, r)
		}
	}
}
