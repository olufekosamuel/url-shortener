// Command server is the process we run. In Go, binaries live under cmd/<name>.
//
// Interview talking point: main.go should be a thin composition root.
// It loads config and wires packages together. It should not contain
// URL validation, ID generation, or SQL.
package main

import (
	"log"

	"github.com/olufekosamuel/url-shortener/internal/config"
	"github.com/olufekosamuel/url-shortener/internal/store"
)

func main() {
	cfg := config.Load()
	_ = store.NewMemoryStore()

	log.Printf("url-shortener skeleton: addr=%s baseURL=%s", cfg.Addr, cfg.BaseURL)
	log.Printf("store: in-memory (lost on restart). next: shortener service, then HTTP")
}
