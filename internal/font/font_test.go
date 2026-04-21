package font

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// standardPath resolves repo root from current test file path
// and returns assets/standard.txt.
func standardPath(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	// thisFile: .../internal/font/font_test.go
	// go up to repo root: internal/font -> internal -> repo root (two levels)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))

	return filepath.Join(repoRoot, "assets", "standard.txt")
}

// TestLoadStandard_SpaceIs8EmptyLines validates that the space glyph exists
// and is represented by 8 blank rows.
func TestLoadStandard_SpaceIs8EmptyLines(t *testing.T) {
	// Load the standard font file into a glyph dictionary.
	f, err := LoadBanner(standardPath(t))
	if err != nil {
		t.Fatalf("LoadStandard returned error: %v", err)
	}

	// The space character (' ') should exist in the font.
	glyph, ok := f.Glyph[' ']
	if !ok {
		t.Fatalf("missing glyph for space ' '")
	}

	// Each glyph must have exactly 8 lines (as per the project spec).
	if len(glyph) != 8 {
		t.Fatalf("space glyph: expected 8 lines, got %d", len(glyph))
	}

	// In the Zone01 standard font, the space glyph is 8 empty lines.
	for i, line := range glyph {
		// TrimRight is assertion-only here; loader behavior preserves spaces.
		if strings.TrimRight(line, " ") != "" {
			t.Fatalf("space glyph: expected blank (spaces allowed) at index %d, got %q", i, line)
		}
	}
}

// TestLoadStandard_ExclamationExistsAndNotAllEmpty validates that a punctuation
// glyph exists and is not entirely blank.
func TestLoadStandard_ExclamationExistsAndNotAllEmpty(t *testing.T) {
	// Load the standard font file into a glyph dictionary.
	f, err := LoadBanner(standardPath(t))
	if err != nil {
		t.Fatalf("LoadStandard returned error: %v", err)
	}

	// The '!' character should exist in the font.
	glyph, ok := f.Glyph['!']
	if !ok {
		t.Fatalf("missing glyph for '!'")
	}

	// Each glyph must have exactly 8 lines.
	if len(glyph) != 8 {
		t.Fatalf("'!' glyph: expected 8 lines, got %d", len(glyph))
	}

	// Ensure the glyph is not entirely empty.
	// NOTE: This trimming is ONLY for the assertion here, not for the loader.
	allEmpty := true
	for _, line := range glyph {
		// TrimRight is assertion-only here; loader behavior preserves spaces.
		if strings.TrimRight(line, " ") != "" {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		t.Fatalf("'!' glyph: expected at least one non-empty line, but all were empty")
	}
}
