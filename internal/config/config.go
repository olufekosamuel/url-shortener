// Package config loads process settings from the environment.
//
// Why a dedicated package: 12-factor style. The same binary can run locally
// and in production by changing ADDR and BASE_URL, not by changing code.
package config

import "os"

type Config struct {
	// Addr is the HTTP listen address, e.g. ":8080".
	Addr string
	// BaseURL is how we present short links, e.g. "http://localhost:8080".
	BaseURL string
}

func Load() Config {
	return Config{
		Addr:    envOr("ADDR", ":8080"),
		BaseURL: envOr("BASE_URL", "http://localhost:8080"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
