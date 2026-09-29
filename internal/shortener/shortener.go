package shortener

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/olufekosamuel/url-shortener/internal/store"
)

var (
	ErrInvalidURL    = errors.New("url must be an absolute http or https URL")
	ErrCodeCollision = errors.New("could not allocate a unique short code")
)

const (
	codeLength  = 7
	maxAttempts = 8
)

// Service validates URLs, generates short codes, and reads/writes the store.
type Service struct {
	store    store.Store
	generate func() (string, error)
}

func New(st store.Store) *Service {
	return &Service{
		store: st,
		generate: func() (string, error) {
			return RandomCode(codeLength)
		},
	}
}

// Shorten validates rawURL, allocates a unique code, and stores the mapping.
func (s *Service) Shorten(ctx context.Context, rawURL string) (string, error) {
	longURL, err := normalizeURL(rawURL)
	if err != nil {
		return "", err
	}

	for i := 0; i < maxAttempts; i++ {
		code, err := s.generate()
		if err != nil {
			return "", err
		}
		err = s.store.Save(ctx, code, longURL)
		if err == nil {
			return code, nil
		}
		if errors.Is(err, store.ErrConflict) {
			continue
		}
		return "", err
	}
	return "", ErrCodeCollision
}

// Resolve returns the long URL for code.
func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", store.ErrNotFound
	}
	return s.store.Get(ctx, code)
}

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrInvalidURL
	}
	if u.Host == "" {
		return "", ErrInvalidURL
	}
	return u.String(), nil
}
