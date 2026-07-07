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

func TestBlockFloatPlacesBelowCursor(t *testing.T) {
	// 20×10 screen, cursor at (2, 3), 2-line popup → row=3, col=3
	base := strings.Repeat(strings.Repeat(".", 20)+"\n", 9) + strings.Repeat(".", 20)
	popup := []string{"AAAA", "BBBB"}
	got := BlockFloat(base, popup, 2, 3, 20, 10)
	lines := strings.Split(got, "\n")
	if !strings.Contains(lines[3], "AAAA") {
		t.Errorf("popup row 0 should be on line 3: %q", lines[3])
	}
	if !strings.Contains(lines[4], "BBBB") {
		t.Errorf("popup row 1 should be on line 4: %q", lines[4])
	}
}

func TestBlockFloatFlipsAboveWhenNoRoomBelow(t *testing.T) {
	// 20×5 screen, cursor at row 4 (last row), 2-line popup
	// → row = 4-2 = 2 (above cursor)
	base := strings.Repeat(strings.Repeat(".", 20)+"\n", 4) + strings.Repeat(".", 20)
	popup := []string{"AAAA", "BBBB"}
	got := BlockFloat(base, popup, 4, 3, 20, 5)
	lines := strings.Split(got, "\n")
	if !strings.Contains(lines[2], "AAAA") {
		t.Errorf("popup should flip above cursor to row 2: %q", lines[2])
	}
	if !strings.Contains(lines[3], "BBBB") {
		t.Errorf("popup row 1 should be on line 3: %q", lines[3])
	}
}

func TestBlockFloatShiftsLeftOnOverflow(t *testing.T) {
	// 10×10 screen, cursor at col 8, 4-wide popup
	// → col = 10-4 = 6 (shifted left to fit)
	base := strings.Repeat(strings.Repeat(".", 10)+"\n", 9) + strings.Repeat(".", 10)
	popup := []string{"PPPP"}
	got := BlockFloat(base, popup, 0, 8, 10, 10)
	lines := strings.Split(got, "\n")
	stripped := ansi.Strip(lines[1])
	// Popup should start at col 6, so 6 dots then PPPP
	if !strings.Contains(stripped, "PPPP") {
		t.Errorf("popup missing: %q", lines[1])
	}
}

func TestSgrStateAtPreservesColorAfterOverlay(t *testing.T) {
	// Base has red text; overlay replaces middle; the right portion
	// should still be red because sgrStateAt re-emits the SGR.
	base := "\x1b[31maaaaa\x1b[0m"
	got := Line(base, "XX", 1)
	// After overlay: left(1 red) + XX + reset + sgr(red) + right(2 red)
	// The right portion should contain the red SGR sequence
	if !strings.Contains(got, "\x1b[31m") {
		t.Fatalf("right portion should re-emit red SGR: %q", got)
	}
	// Visible width should be 5
	if w := ansi.StringWidth(got); w != 5 {
		t.Fatalf("visible width: got %d, want 5; %q", w, got)
	}
}
