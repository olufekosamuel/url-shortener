package handler

// Handler owns HTTP: routing, JSON, status codes, redirects.
//
// Planned routes:
//   POST /v1/urls   create a short link
//   GET  /{code}    302/301 to the original URL
//
// This package stays thin. It calls shortener.Service and does not import
// a database driver.
type Handler struct {
	// Service is injected from main.go.
}
