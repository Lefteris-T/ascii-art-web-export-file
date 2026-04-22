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

func TestAsciiArt_WrongMethodStillWorks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rr := httptest.NewRecorder()

	asciiArt(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

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
