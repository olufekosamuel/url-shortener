// Command server is the process we run. In Go, binaries live under cmd/<name>.
//
// Interview talking point: main.go should be a thin composition root.
// It loads config and wires packages together. It should not contain
// URL validation, ID generation, or SQL.
package main

import (
	"log"

	"github.com/samuelolufeko/url-shortener/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("url-shortener skeleton listening config: addr=%s baseURL=%s", cfg.Addr, cfg.BaseURL)
	log.Printf("next: implement store, then shortener service, then HTTP handlers")
}
