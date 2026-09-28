package handler

// Handler will own HTTP: routing, JSON, status codes, redirects.
//
// Planned routes (not implemented yet):
//   POST /v1/urls   create a short link
//   GET  /{code}    302/301 to the original URL
//
// This package should stay thin. It calls shortener.Service.
// It must not import a database driver. If it needs storage, it is wired
// through the service, which talks to store.Store.
type Handler struct {
	// Service will be injected from main.go later.
}
