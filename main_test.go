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

func testMux() *http.ServeMux {
	mux := http.NewServeMux()
	handlers.Register(mux)
	return mux
}

func postForm(t *testing.T, target string, values url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	testMux().ServeHTTP(rec, req)
	return rec
}

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

func TestHomeHandler_MethodNotAllowedAsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHomeHandler_UnknownRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

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

func TestASCIIArtHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rec := httptest.NewRecorder()

	testMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

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
