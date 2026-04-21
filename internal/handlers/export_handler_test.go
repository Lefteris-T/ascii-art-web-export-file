package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestExportHandler_Success(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	form := url.Values{}
	form.Set("text", "A")
	form.Set("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("expected Content-Type to contain %q, got %q", "text/plain", contentType)
	}

	contentDisposition := rr.Header().Get("Content-Disposition")
	if contentDisposition == "" {
		t.Fatal("expected Content-Disposition header to be set")
	}

	if !strings.Contains(contentDisposition, "attachment") {
		t.Fatalf("expected Content-Disposition to contain %q, got %q", "attachment", contentDisposition)
	}

	if !strings.Contains(contentDisposition, ".txt") {
		t.Fatalf("expected Content-Disposition to contain .txt filename, got %q", contentDisposition)
	}

	if rr.Body.Len() == 0 {
		t.Fatal("expected non-empty export body")
	}

	if strings.Contains(strings.ToLower(rr.Body.String()), "<html") {
		t.Fatal("expected plain text response body, got HTML")
	}
}

func TestExportHandler_WrongMethod(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/export", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestExportHandler_InvalidBanner(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "invalid")

	req := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestExportHandler_RenderError(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	form := url.Values{}
	form.Set("text", string(rune(1)))
	form.Set("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
