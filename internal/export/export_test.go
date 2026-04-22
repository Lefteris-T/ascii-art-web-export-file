package export

import "testing"

// TestBuild_EmptyContent verifies that exporting an empty render result does
// not add placeholder text or newlines.
func TestBuild_EmptyContent(t *testing.T) {
	got := Build("")
	want := ""

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// TestBuild_PassThroughSingleLine proves the export layer does not transform a
// simple one-line result.
func TestBuild_PassThroughSingleLine(t *testing.T) {
	input := "A"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

// TestBuild_PreservesMultipleLines protects line breaks in rendered ASCII art.
func TestBuild_PreservesMultipleLines(t *testing.T) {
	input := "ABCD\nEFGH\nIJKL\n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

// TestBuild_PreservesLeadingAndTrailingSpaces protects visible ASCII layout,
// where spaces at either side of a line can be meaningful.
func TestBuild_PreservesLeadingAndTrailingSpaces(t *testing.T) {
	input := "  ABC\nDEF  \n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

// TestBuild_DoesNotRemoveEmptyLines ensures vertical spacing survives export.
func TestBuild_DoesNotRemoveEmptyLines(t *testing.T) {
	input := "AAAA\n\nBBBB\n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

// TestBuild_ByteForByteIdentity is the broad guard that export content remains
// exactly the same string that the render flow produced.
func TestBuild_ByteForByteIdentity(t *testing.T) {
	input := " \nASCII ART\n\nline with spaces  \n\tindent\n"
	got := Build(input)

	if len(got) != len(input) {
		t.Fatalf("expected length %d, got %d", len(input), len(got))
	}

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}
