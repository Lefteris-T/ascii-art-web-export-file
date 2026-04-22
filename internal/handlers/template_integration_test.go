package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomePage_NoExportFormWithoutResult(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	home(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()

	if strings.Contains(body, `action="/export"`) {
		t.Fatalf("did not expect export form when HasResult is false")
	}
}

func TestAsciiArtPage_ShowsExportFormWhenResultExists(t *testing.T) {
	form := strings.NewReader("text=hello&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	asciiArt(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()

	if !strings.Contains(body, `action="/export"`) {
		t.Fatalf("expected export form action to be present")
	}

	if !strings.Contains(body, `method="POST"`) {
		t.Fatalf("expected export form method POST to be present")
	}

	if !strings.Contains(body, `type="hidden" name="text"`) {
		t.Fatalf("expected hidden text field to be present")
	}

	if !strings.Contains(body, `type="hidden" name="banner"`) {
		t.Fatalf("expected hidden banner field to be present")
	}

	if !strings.Contains(body, `value="standard"`) {
		t.Fatalf("expected hidden banner value to be present")
	}
}
