package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestExportHandler_Success verifies the happy path for downloading generated
// ASCII art as a plain text attachment.
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

	// Successful export must identify the body as text, not HTML.
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("expected Content-Type to contain %q, got %q", "text/plain", contentType)
	}

	// Content-Disposition is what tells the browser to download a .txt file.
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

// TestExportHandler_WrongMethod protects the route contract: export accepts
// POST only and rejects direct GET requests.
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

// TestExportHandler_InvalidBanner ensures export uses the same banner
// validation rules as normal rendering.
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

// TestExportHandler_RenderError checks that invalid render input is surfaced as
// a bad request instead of producing a downloadable file.
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

// extractASCIIResult pulls the rendered ASCII block out of the HTML response so
// the integration test can compare it directly with the exported text body.
func extractASCIIResult(t *testing.T, html string) string {
	t.Helper()

	startTag := `<pre id="ascii-result">`
	endTag := `</pre>`

	start := strings.Index(html, startTag)
	if start == -1 {
		t.Fatal("could not find opening ascii result tag in rendered HTML")
	}
	start += len(startTag)

	end := strings.Index(html[start:], endTag)
	if end == -1 {
		t.Fatal("could not find closing ascii result tag in rendered HTML")
	}

	return html[start : start+end]
}

// TestExportMatchesRenderedOutput proves export and page rendering use the same
// server-side source of truth for ASCII generation.
func TestExportMatchesRenderedOutput(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	form := url.Values{}
	form.Set("text", "A")
	form.Set("banner", "standard")

	renderReq := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	renderReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	renderRec := httptest.NewRecorder()

	mux.ServeHTTP(renderRec, renderReq)

	if renderRec.Code != http.StatusOK {
		t.Fatalf("render request: expected status %d, got %d", http.StatusOK, renderRec.Code)
	}

	renderedHTML := renderRec.Body.String()
	renderedASCII := extractASCIIResult(t, renderedHTML)

	exportReq := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader(form.Encode()))
	exportReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	exportRec := httptest.NewRecorder()

	mux.ServeHTTP(exportRec, exportReq)

	if exportRec.Code != http.StatusOK {
		t.Fatalf("export request: expected status %d, got %d", http.StatusOK, exportRec.Code)
	}

	exportedASCII := exportRec.Body.String()

	if exportedASCII != renderedASCII {
		t.Fatalf("expected export output to match rendered output exactly\n--- rendered ---\n%q\n--- exported ---\n%q", renderedASCII, exportedASCII)
	}
}
