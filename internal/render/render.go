package render

import (
	"ascii-art-web-export-file/internal/font"
	"fmt"
	"strings"
)

// Render converts plain text into ASCII-art using the provided font glyph map.
// It supports multiline input via newline characters and preserves output layout exactly.
func Render(text string, f font.Font) (string, error) {
	if text == "" {
		return "", nil
	}

	// Normalize Windows line endings just in case.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// If the input is only newline characters, preserve exactly that count.
	// This avoids an extra blank line from Split's trailing empty element.
	if strings.Trim(text, "\n") == "" {
		return strings.Repeat("\n", strings.Count(text, "\n")), nil
	}

	lines := strings.Split(text, "\n")

	var b strings.Builder

	for _, line := range lines {
		// Empty input line corresponds to a single blank output line.
		if line == "" {
			b.WriteByte('\n')
			continue
		}

		// Render non-empty input as an 8-row ASCII block.
		for row := 0; row < 8; row++ {
			for _, ch := range line {
				glyph, ok := f.Glyph[ch]
				if !ok {
					return "", fmt.Errorf("missing glyph for %q", ch)
				}
				if row >= len(glyph) {
					return "", fmt.Errorf("invalid glyph data for %q (row %d out of range)", ch, row)
				}
				b.WriteString(glyph[row])
			}
			b.WriteByte('\n')
		}
	}

	return b.String(), nil
}
