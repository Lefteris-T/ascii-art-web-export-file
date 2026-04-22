package export

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestBuild_TXTEmptyContent verifies that exporting an empty render result does
// not add placeholder text or newlines.
func TestBuild_TXTEmptyContent(t *testing.T) {
	content, contentType, filename, err := Build("txt", "", "standard", "")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if content != "" {
		t.Fatalf("expected empty content, got %q", content)
	}
	if contentType != "text/plain; charset=utf-8" {
		t.Fatalf("expected text content type, got %q", contentType)
	}
	if filename != "ascii-art.txt" {
		t.Fatalf("expected txt filename, got %q", filename)
	}
}

// TestBuild_TXTPreservesByteForByteIdentity proves the text export does not
// transform the rendered ASCII result.
func TestBuild_TXTPreservesByteForByteIdentity(t *testing.T) {
	input := " \nASCII ART\n\nline with spaces  \n\tindent\n"

	content, _, _, err := Build("txt", "source text", "standard", input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if len(content) != len(input) {
		t.Fatalf("expected length %d, got %d", len(input), len(content))
	}
	if content != input {
		t.Fatalf("expected %q, got %q", input, content)
	}
}

// TestBuild_HTMLWrapsEscapedResult checks that HTML export wraps the result in
// a document and escapes the ASCII content before placing it in a pre block.
func TestBuild_HTMLWrapsEscapedResult(t *testing.T) {
	content, contentType, filename, err := Build("html", "source", "standard", "<ASCII & ART>")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected html content type, got %q", contentType)
	}
	if filename != "ascii-art.html" {
		t.Fatalf("expected html filename, got %q", filename)
	}
	if !strings.Contains(content, "<!DOCTYPE html>") {
		t.Fatalf("expected html document, got %q", content)
	}
	if !strings.Contains(content, "&lt;ASCII &amp; ART&gt;") {
		t.Fatalf("expected escaped result in html export, got %q", content)
	}
}

// TestBuild_JSONIncludesMetadata verifies that JSON export includes the source
// text, selected banner, and rendered result.
func TestBuild_JSONIncludesMetadata(t *testing.T) {
	content, contentType, filename, err := Build("json", "Hello", "standard", "ASCII\n")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if contentType != "application/json" {
		t.Fatalf("expected json content type, got %q", contentType)
	}
	if filename != "ascii-art.json" {
		t.Fatalf("expected json filename, got %q", filename)
	}

	var payload Payload
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		t.Fatalf("json export did not unmarshal: %v", err)
	}

	if payload.Text != "Hello" {
		t.Fatalf("expected text %q, got %q", "Hello", payload.Text)
	}
	if payload.Banner != "standard" {
		t.Fatalf("expected banner %q, got %q", "standard", payload.Banner)
	}
	if payload.Result != "ASCII\n" {
		t.Fatalf("expected result %q, got %q", "ASCII\n", payload.Result)
	}
}

// TestBuild_UnsupportedFormat verifies unknown formats are rejected instead of
// silently falling back to another file type.
func TestBuild_UnsupportedFormat(t *testing.T) {
	content, contentType, filename, err := Build("pdf", "Hello", "standard", "ASCII")
	if err == nil {
		t.Fatal("expected unsupported format error")
	}
	if content != "" || contentType != "" || filename != "" {
		t.Fatalf("expected empty export values on error, got content=%q contentType=%q filename=%q", content, contentType, filename)
	}
}
