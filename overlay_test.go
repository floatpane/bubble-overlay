package overlay

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestLinePaintsOverlayInMiddleOfPlainBase(t *testing.T) {
	got := Line("aaaaaaaaaa", "XXX", 3)
	if got != "aaaXXX\x1b[0maaaa" {
		t.Fatalf("got %q", got)
	}
}

func TestLineAtColumnZeroReplacesPrefix(t *testing.T) {
	got := Line("aaaaa", "XX", 0)
	if got != "XX\x1b[0maaa" {
		t.Fatalf("got %q", got)
	}
}

func TestLinePadsBeyondVisibleWidth(t *testing.T) {
	got := Line("ab", "X", 5)
	if got != "ab   X\x1b[0m" {
		t.Fatalf("got %q", got)
	}
}

func TestLineEmptyOverlayReturnsBaseUnchanged(t *testing.T) {
	if got := Line("abc", "", 1); got != "abc" {
		t.Fatalf("got %q, want %q", got, "abc")
	}
}

func TestLinePreservesSurroundingAnsi(t *testing.T) {
	base := "\x1b[31maaaaa\x1b[0m" // red "aaaaa"
	got := Line(base, "XX", 1)
	// Must still contain the overlay and a reset after it.
	if !strings.Contains(got, "XX\x1b[0m") {
		t.Fatalf("missing overlay + reset in %q", got)
	}
	// Visible width must equal max(baseWidth, col + overlayWidth) — 5 here.
	if w := ansi.StringWidth(got); w != 5 {
		t.Fatalf("got visible width %d, want 5; output=%q", w, got)
	}
}

func TestBlockMultilineGrowsBase(t *testing.T) {
	base := "row0\nrow1"
	got := Block(base, []string{"AA", "BB", "CC"}, 1, 1)
	lines := strings.Split(got, "\n")
	// Overlay starts at row 1 and is 3 lines tall, so rows 2 and 3 are
	// appended past the original 2-line base.
	if len(lines) != 4 {
		t.Fatalf("got %d lines, want 4: %q", len(lines), lines)
	}
	if lines[0] != "row0" {
		t.Fatalf("row 0 mutated: %q", lines[0])
	}
	if !strings.Contains(lines[1], "AA") {
		t.Fatalf("line 1 = %q, want overlay AA", lines[1])
	}
	if !strings.Contains(lines[2], "BB") {
		t.Fatalf("line 2 = %q, want overlay BB", lines[2])
	}
	if !strings.Contains(lines[3], "CC") {
		t.Fatalf("line 3 = %q, want overlay CC", lines[3])
	}
}

func TestBlockEmptyOverlayReturnsBase(t *testing.T) {
	base := "hello\nworld"
	if got := Block(base, nil, 0, 0); got != base {
		t.Fatalf("got %q, want unchanged", got)
	}
}
