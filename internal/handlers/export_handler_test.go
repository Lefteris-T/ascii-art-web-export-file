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
