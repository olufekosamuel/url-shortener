// Command server is the process entrypoint.
// main loads config, constructs dependencies, and runs the HTTP server.
// URL validation, code generation, and storage live in internal packages.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/olufekosamuel/url-shortener/internal/config"
	"github.com/olufekosamuel/url-shortener/internal/handler"
	"github.com/olufekosamuel/url-shortener/internal/shortener"
	"github.com/olufekosamuel/url-shortener/internal/store"
)

func main() {
	cfg := config.Load()
	svc := shortener.New(store.NewMemoryStore())
	h := handler.New(svc, cfg.BaseURL)

	// Timeouts stop slow or idle clients from holding connections forever.
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s (baseURL=%s, store=memory)", cfg.Addr, cfg.BaseURL)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	case <-ctx.Done():
		log.Printf("shutting down")
		// Let in-flight requests finish, but not forever.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown: %v", err)
		}
	}
}
