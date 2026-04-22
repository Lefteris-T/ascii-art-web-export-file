package render

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"ascii-art-web-export-file/internal/font"
)

// repoRoot returns the repository root directory based on the location of this test file.
func repoRoot(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	// thisFile: .../internal/render/render_test.go
	// up to repo root: internal/render -> internal -> repo root (two levels up from internal)
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

func loadStandardFont(t *testing.T) font.Font {
	t.Helper()

	f, err := font.LoadBanner(filepath.Join(repoRoot(t), "assets", "standard.txt"))
	if err != nil {
		t.Fatalf("LoadBanner: %v", err)
	}

	return f
}

func expectedRenderedLine(t *testing.T, line string, f font.Font) string {
	t.Helper()

	var b strings.Builder
	for row := 0; row < 8; row++ {
		for _, ch := range line {
			glyph, ok := f.Glyph[ch]
			if !ok {
				t.Fatalf("missing test glyph for %q", ch)
			}
			if row >= len(glyph) {
				t.Fatalf("test glyph for %q has %d rows", ch, len(glyph))
			}
			b.WriteString(glyph[row])
		}
		b.WriteByte('\n')
	}

	return b.String()
}

func expectedRenderedText(t *testing.T, text string, f font.Font) string {
	t.Helper()

	if text == "" {
		return ""
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	if strings.Trim(text, "\n") == "" {
		return strings.Repeat("\n", strings.Count(text, "\n"))
	}

	var b strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			b.WriteByte('\n')
			continue
		}
		b.WriteString(expectedRenderedLine(t, line, f))
	}

	return b.String()
}

// TestRender_Empty verifies empty input returns empty output.
func TestRender_Empty(t *testing.T) {
	f := loadStandardFont(t)

	out, err := Render("", f)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != "" {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", "", out)
	}
}

// TestRender_A verifies rendering of a single character glyph.
func TestRender_A(t *testing.T) {
	f := loadStandardFont(t)
	input := "A"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_Hello verifies rendering of a simple lowercase word.
func TestRender_Hello(t *testing.T) {
	f := loadStandardFont(t)
	input := "Hello"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_HelloThere_Newline verifies escaped newline behavior across two rendered blocks.
func TestRender_HelloThere_Newline(t *testing.T) {
	f := loadStandardFont(t)
	input := "Hello\nThere"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_Newline_only verifies newline-only input preserves line count correctly.
func TestRender_Newline_only(t *testing.T) {
	f := loadStandardFont(t)
	input := "\n"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_Hello_TrailingNewline verifies output when rendered text ends with a newline.
func TestRender_Hello_TrailingNewline(t *testing.T) {
	f := loadStandardFont(t)
	input := "Hello\n"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_DoubleNewline verifies two consecutive newline characters are preserved.
func TestRender_DoubleNewline(t *testing.T) {
	f := loadStandardFont(t)
	input := "\n\n"
	expected := expectedRenderedText(t, input, f)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}
