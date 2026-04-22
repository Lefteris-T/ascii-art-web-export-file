package export

import "testing"

func TestBuild_EmptyContent(t *testing.T) {
	got := Build("")
	want := ""

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuild_PassThroughSingleLine(t *testing.T) {
	input := "A"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestBuild_PreservesMultipleLines(t *testing.T) {
	input := "ABCD\nEFGH\nIJKL\n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestBuild_PreservesLeadingAndTrailingSpaces(t *testing.T) {
	input := "  ABC\nDEF  \n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestBuild_DoesNotRemoveEmptyLines(t *testing.T) {
	input := "AAAA\n\nBBBB\n"
	got := Build(input)

	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

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
