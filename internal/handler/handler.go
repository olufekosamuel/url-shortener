package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/olufekosamuel/url-shortener/internal/shortener"
	"github.com/olufekosamuel/url-shortener/internal/store"
)

// maxBodyBytes caps the create request body so a client cannot stream
// an unbounded payload into the JSON decoder.
const maxBodyBytes = 1 << 20

// Handler owns HTTP: routing, JSON, status codes, redirects.
//
// Routes:
//
//	POST /v1/urls   create a short link
//	GET  /{code}    302 to the original URL
//
// This package stays thin. It calls shortener.Service and does not import
// a database driver.
type Handler struct {
	svc     *shortener.Service
	baseURL string
}

func New(svc *shortener.Service, baseURL string) *Handler {
	return &Handler{
		svc:     svc,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// Routes returns the router. Method and wildcard patterns need Go 1.22+.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/urls", h.create)
	mux.HandleFunc("GET /{code}", h.redirect)
	return mux
}

type createRequest struct {
	URL string `json:"url"`
}

type createResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	code, err := h.svc.Shorten(r.Context(), req.URL)
	switch {
	case err == nil:
	case errors.Is(err, shortener.ErrInvalidURL):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	default:
		log.Printf("shorten %q: %v", req.URL, err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusCreated, createResponse{
		Code:     code,
		ShortURL: h.baseURL + "/" + code,
	})
}

// redirect uses 302, not 301: browsers cache 301s indefinitely, which would
// hide later changes to a link and skip any future click counting.
func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	longURL, err := h.svc.Resolve(r.Context(), r.PathValue("code"))
	switch {
	case err == nil:
		http.Redirect(w, r, longURL, http.StatusFound)
	case errors.Is(err, store.ErrNotFound):
		http.NotFound(w, r)
	default:
		log.Printf("resolve %q: %v", r.PathValue("code"), err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
