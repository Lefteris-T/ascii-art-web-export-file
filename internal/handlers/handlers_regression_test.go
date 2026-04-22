package handlers

import (
	"ascii-art-web-export-file/internal/font"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestHome_GetStillWorks guards the existing home page while export work is
// added around it.
func TestHome_GetStillWorks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	home(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "ASCII Art Generator") {
		t.Fatalf("expected home page content to be present")
	}
}

// TestAsciiArt_PostStillWorks confirms the original render submission path
// still returns a result page.
func TestAsciiArt_PostStillWorks(t *testing.T) {
	form := url.Values{}
	form.Set("text", "A")
	form.Set("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	asciiArt(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "ASCII Art Result") {
		t.Fatalf("expected result section to be present")
	}
}

// TestAsciiArt_WrongMethodStillWorks preserves existing invalid method
// behavior for the render endpoint.
func TestAsciiArt_WrongMethodStillWorks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rr := httptest.NewRecorder()

	asciiArt(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestAsciiArt_InvalidBannerStillWorks ensures invalid banner handling did not
// change while adding export.
func TestAsciiArt_InvalidBannerStillWorks(t *testing.T) {
	form := url.Values{}
	form.Set("text", "A")
	form.Set("banner", "invalid-banner")

	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	asciiArt(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAsciiArt_MissingAssetStillWorks(t *testing.T) {
	form := url.Values{}
	form.Set("text", "A")
	form.Set("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	originalLoadBanner := loadBanner
	// Stub banner loading so the handler follows the same path as a missing
	// asset on disk without renaming real project files.
	loadBanner = func(path string) (font.Font, error) {
		return font.Font{}, errors.New("missing asset stub")
	}
	defer func() {
		loadBanner = originalLoadBanner
	}()

	asciiArt(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

// TestRenderErrorPage_StillWorks checks that unknown routes still render the
// shared error page after route additions.
func TestRenderErrorPage_StillWorks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rr := httptest.NewRecorder()

	home(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Not Found") {
		t.Fatalf("expected not found page content")
	}
}
