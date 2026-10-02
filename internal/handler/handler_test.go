package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/olufekosamuel/url-shortener/internal/shortener"
	"github.com/olufekosamuel/url-shortener/internal/store"
)

func newTestServer() http.Handler {
	svc := shortener.New(store.NewMemoryStore())
	return New(svc, "http://short.test/").Routes()
}

func post(t *testing.T, h http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/urls", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateAndRedirect(t *testing.T) {
	h := newTestServer()

	rec := post(t, h, `{"url":"https://example.com/a"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201; body %s", rec.Code, rec.Body)
	}
	var resp createResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Code == "" {
		t.Fatal("response code is empty")
	}
	if want := "http://short.test/" + resp.Code; resp.ShortURL != want {
		t.Fatalf("short_url = %q, want %q", resp.ShortURL, want)
	}

	req := httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("GET status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/a" {
		t.Fatalf("Location = %q, want https://example.com/a", loc)
	}
}

func TestCreateBadRequests(t *testing.T) {
	h := newTestServer()

	for _, body := range []string{
		``,
		`not json`,
		`{"url":"example.com"}`,
		`{"url":"ftp://example.com"}`,
		`{"link":"https://example.com"}`,
	} {
		rec := post(t, h, body)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("POST %q: status = %d, want 400", body, rec.Code)
		}
	}
}

func TestRedirectUnknownCode(t *testing.T) {
	h := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestWrongMethod(t *testing.T) {
	h := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/v1/urls", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	// GET /v1/urls does not match GET /{code} (single segment), and
	// POST /v1/urls only allows POST.
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
