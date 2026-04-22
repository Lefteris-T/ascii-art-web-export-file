package export

// Build returns the final ASCII content that should be written to the exported
// file. It intentionally leaves the content unchanged so spacing, newlines, and
// empty lines are preserved byte-for-byte.
func Build(content string) string {
	return content
}
