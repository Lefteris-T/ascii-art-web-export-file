package main

import (
	"ascii-art-web-export-file/internal/font"
	"ascii-art-web-export-file/internal/handlers"
	"ascii-art-web-export-file/internal/render"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

// testMux builds the same application routes used by main without starting a
// real HTTP server.
func testMux() *http.ServeMux {
	mux := http.NewServeMux()
	handlers.Register(mux)
	return mux
}

// postForm submits URL-encoded form data through the test mux and returns the
// recorded response for assertions.
func postForm(t *testing.T, target string, values url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	testMux().ServeHTTP(rec, req)
	return rec
}

// TestHomeHandler_OK verifies that the root route serves the generator page.
func TestHomeHandler_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "ASCII Art Generator") {
		t.Fatalf("response body missing page title: %q", body)
	}
}

// TestHomeHandler_MethodNotAllowedAsBadRequest preserves the project contract
// that unsupported methods return 400 instead of the default 405.
func TestHomeHandler_MethodNotAllowedAsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestHomeHandler_UnknownRoute checks that the root handler rejects paths that
// do not belong to the application.
func TestHomeHandler_UnknownRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// TestASCIIArtHandler_OK compares the rendered page against output generated
// directly by the font and render packages.
func TestASCIIArtHandler_OK(t *testing.T) {
	values := url.Values{
		"text":   {"Hello"},
		"banner": {"standard"},
	}

	rec := postForm(t, "/ascii-art", values)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	f, err := font.LoadBanner(filepath.Join("assets", "standard.txt"))
	if err != nil {
		t.Fatalf("LoadBanner: %v", err)
	}

	expected, err := render.Render("Hello", f)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "ASCII Art Result") {
		t.Fatalf("response body missing result page title: %q", body)
	}
	if !strings.Contains(body, html.EscapeString(expected)) {
		t.Fatalf("response body missing rendered output")
	}
}

// TestASCIIArtHandler_InvalidMethod verifies that the render endpoint only
// accepts POST form submissions.
func TestASCIIArtHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestASCIIArtHandler_InvalidBanner verifies that unsupported banner names are
// rejected before any render output is produced.
func TestASCIIArtHandler_InvalidBanner(t *testing.T) {
	values := url.Values{
		"text":   {"Hello"},
		"banner": {"invalid"},
	}

	rec := postForm(t, "/ascii-art", values)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestASCIIArtHandler_InvalidInput covers characters that are not available in
// the printable ASCII banner glyph map.
func TestASCIIArtHandler_InvalidInput(t *testing.T) {
	values := url.Values{
		"text":   {"Hello\tWorld"},
		"banner": {"standard"},
	}

	rec := postForm(t, "/ascii-art", values)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestASCIIArtHandler_EscapedNewline confirms that form text containing the
// literal sequence "\n" is normalized before rendering.
func TestASCIIArtHandler_EscapedNewline(t *testing.T) {
	values := url.Values{
		"text":   {"Hello\\nWorld"},
		"banner": {"standard"},
	}

	rec := postForm(t, "/ascii-art", values)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	f, err := font.LoadBanner(filepath.Join("assets", "standard.txt"))
	if err != nil {
		t.Fatalf("LoadBanner: %v", err)
	}

	expected, err := render.Render("Hello\nWorld", f)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if !strings.Contains(rec.Body.String(), html.EscapeString(expected)) {
		t.Fatalf("response body missing multiline rendered output")
	}
}
