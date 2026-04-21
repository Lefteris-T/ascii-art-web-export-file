package font

import (
	"fmt"
	"os"
	"strings"
)

// Font stores the parsed ASCII-art glyphs, keyed by printable rune.
// Each glyph is expected to contain exactly 8 rows.
type Font struct {
	Glyph map[rune][]string
}

// LoadBanner parses the banner file into a glyph map.
// It expects printable ASCII characters in order from 32 to 126, where each
// character is represented by 8 rows separated by one newline separator.
func LoadBanner(path string) (Font, error) {
	lines, err := readLines(path)
	if err != nil {
		return Font{}, err
	}

	// The banner contains 95 printable ASCII characters (32..126).
	// File layout is: 1 header line + (95 glyphs * 8 rows) + 94 separator lines.
	const needed = 1 + 95*8 + 94 // 855 total lines after readLines normalization

	if len(lines) < needed {
		return Font{}, fmt.Errorf("banner file too short: got %d lines, need at least %d", len(lines), needed)
	}

	glyphs := make(map[rune][]string)
	// The first line is a header/offset line in standard.txt.
	i := 1
	for ascii := 32; ascii <= 126; ascii++ {
		// Guard: we must have 8 contiguous lines available for each glyph.
		if i+8 > len(lines) {
			return Font{}, fmt.Errorf("unexpected EOF while parsing glyph %d", ascii)
		}
		glyphs[rune(ascii)] = lines[i : i+8]
		i += 8

		// Banner format uses one separator line between glyphs.
		// The last glyph has no separator to skip.
		if ascii != 126 {
			if i >= len(lines) {
				return Font{}, fmt.Errorf("unexpected EOF while skipping separator after glyph %d", ascii)
			}
			i++
		}
	}
	return Font{Glyph: glyphs}, nil
}

// readLines reads banner content while preserving spaces exactly as stored.
// It normalizes CRLF to LF and removes only the trailing split artifact caused
// by a final newline at EOF.
func readLines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := string(b)
	s = strings.ReplaceAll(s, "\r\n", "\n")

	lines := strings.Split(s, "\n")

	// Remove ONLY the extra empty element created by a trailing newline.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}
