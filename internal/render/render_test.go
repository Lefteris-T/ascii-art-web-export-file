package render

import (
	"os"
	"path/filepath"
	"runtime"
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

// readFixture loads a test input/expected fixture file as a raw string.
func readFixture(t *testing.T, relPath string) string {
	t.Helper()

	b, err := os.ReadFile(relPath)
	if err != nil {
		t.Fatalf("read fixture %q: %v", relPath, err)
	}
	return string(b)
}

// TestRender_Empty verifies empty input returns empty output.
func TestRender_Empty(t *testing.T) {
	// Load the standard font once for the test.
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}

	// Read input/expected fixtures.
	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "empty.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "empty.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

	// Call Render and compare byte-for-byte.
	out, err := Render(input, f)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}

// TestRender_A verifies rendering of a single character glyph.
func TestRender_A(t *testing.T) {
	// Load the standard font file into a glyph dictionary.
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}

	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "A.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "A.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

	// Call Render and compare byte-for-byte.
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
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}
	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

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
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}
	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello_there_newline.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello_there_newline.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

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
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}
	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "newline_only.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "newline_only.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

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
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}

	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello_trailing_newline.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "hello_trailing_newline.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

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
	stdPath := filepath.Join(repoRoot(t), "assets", "standard.txt")
	f, err := font.LoadBanner(stdPath)
	if err != nil {
		t.Fatalf("LoadStandard error: %v", err)
	}

	inPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "double_newline.input.txt")
	expPath := filepath.Join(repoRoot(t), "internal", "render", "testdata", "double_newline.expected.txt")

	input := readFixture(t, inPath)
	expected := readFixture(t, expPath)

	out, err := Render(input, f)

	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if out != expected {
		t.Fatalf("output mismatch\n--- expected ---\n%q\n--- got ---\n%q", expected, out)
	}
}
