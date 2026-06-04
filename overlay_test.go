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

func TestCenterPlacesPopupInMiddle(t *testing.T) {
	// 10×5 screen, 4×1 popup → col=3, row=2
	base := strings.Repeat(strings.Repeat(".", 10)+"\n", 4) + strings.Repeat(".", 10)
	popup := "PPPP"
	got := Center(base, popup, 10, 5)
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d: %q", len(lines), got)
	}
	if !strings.Contains(lines[2], "PPPP") {
		t.Errorf("popup not on row 2: %q", lines[2])
	}
	// Cells to the left and right of the popup must survive.
	if !strings.HasPrefix(ansi.Strip(lines[2]), "...") {
		t.Errorf("left background missing on popup row: %q", lines[2])
	}
	if !strings.HasSuffix(ansi.Strip(strings.TrimRight(lines[2], " ")), "...") {
		t.Errorf("right background missing on popup row: %q", lines[2])
	}
}

func TestCenterMultilinePopup(t *testing.T) {
	// 20×10 screen, 4×3 popup → col=8, row=3 (floor division)
	base := strings.Repeat(strings.Repeat(".", 20)+"\n", 9) + strings.Repeat(".", 20)
	popup := "AAAA\nBBBB\nCCCC"
	got := Center(base, popup, 20, 10)
	lines := strings.Split(got, "\n")
	for i, want := range []string{"AAAA", "BBBB", "CCCC"} {
		if !strings.Contains(lines[3+i], want) {
			t.Errorf("row %d: want %q in %q", 3+i, want, lines[3+i])
		}
	}
}

func TestCenterClampsToZeroWhenPopupLargerThanScreen(t *testing.T) {
	// popup wider and taller than screen — should not panic, row/col clamped to 0
	base := "ab\ncd"
	popup := "XXXXXXXXXX\nYYYYYYYYYY\nZZZZZZZZZZ"
	got := Center(base, popup, 4, 2)
	if !strings.Contains(got, "XXXXXXXXXX") {
		t.Errorf("popup content missing: %q", got)
	}
}
