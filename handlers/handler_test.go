package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"url-shortener/storage"
)

// TestShortenThenRedirect_PreservesFragment is the end-to-end guard: a fragment
// must survive the round trip and reach the client in the Location header.
func TestShortenThenRedirect_PreservesFragment(t *testing.T) {
	const longURL = "https://app.example.com/#/dashboard"

	h := NewHandler(storage.NewStore())

	form := url.Values{"url": {longURL}}
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.ShortenURL(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /shorten = %d, want %d (body %q)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	code := strings.TrimPrefix(strings.TrimSpace(rec.Body.String()), "http://shorty.url/")

	req = httptest.NewRequest("GET", "/"+code, nil)
	rec = httptest.NewRecorder()
	h.RedirectURL(rec, req)

	if got := rec.Header().Get("Location"); got != longURL {
		t.Errorf("Location header\n got %q\nwant %q", got, longURL)
	}
}
