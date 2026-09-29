// Command server is the process entrypoint.
// main loads config and constructs dependencies. URL validation,
// code generation, and storage live in internal packages.
package main

import (
	"log"

	"github.com/olufekosamuel/url-shortener/internal/config"
	"github.com/olufekosamuel/url-shortener/internal/shortener"
	"github.com/olufekosamuel/url-shortener/internal/store"
)

func main() {
	cfg := config.Load()
	_ = shortener.New(store.NewMemoryStore())

	log.Printf("addr=%s baseURL=%s", cfg.Addr, cfg.BaseURL)
	log.Printf("store: memory; shortener ready; HTTP not started")
}
