// Package overlay paints rectangular blocks of styled text on top of an
// existing ANSI-styled string at a given (row, col) cell position.
//
// It's designed for Bubble Tea / lipgloss views where you need to render a
// modal, popup, tooltip, or floating panel over a fully-styled base view
// without manually unwrapping and re-emitting SGR sequences. Existing
// styling around the overlaid region is preserved; cells underneath are
// replaced.
package overlay

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Block paints the lines of block on top of base starting at the
// (row, col) cell position. Lines that extend past the bottom of base
// are appended. The result preserves existing ANSI styling around the
// overlaid region.
//
// Both row and col are 0-indexed. col is measured in terminal cells
// (not bytes or runes), so wide-character handling matches what the
// terminal will actually render.
func Block(base string, block []string, row, col int) string {
	if len(block) == 0 {
		return base
	}
	lines := strings.Split(base, "\n")
	for i, overlay := range block {
		r := row + i
		for r >= len(lines) {
			lines = append(lines, "")
		}
		lines[r] = Line(lines[r], overlay, col)
	}
	return strings.Join(lines, "\n")
}

// Line returns base with overlay painted starting at column col.
// Existing cells under the overlay are removed; cells to the left and
// right of the overlay are preserved with their ANSI styling intact.
// When col exceeds the visible width of base the gap is padded with
// spaces.
func Line(base, overlay string, col int) string {
	if overlay == "" {
		return base
	}
	overlayWidth := ansi.StringWidth(overlay)
	baseWidth := ansi.StringWidth(base)

	left := ansi.Truncate(base, col, "")
	leftWidth := ansi.StringWidth(left)

	var pad string
	if leftWidth < col {
		pad = strings.Repeat(" ", col-leftWidth)
	}

	var right string
	rightStart := col + overlayWidth
	if rightStart < baseWidth {
		right = ansi.Cut(base, rightStart, baseWidth)
	}

	// Reset SGR after the overlay so the overlay's styles don't bleed
	// into the surrounding cells (the rest of the row may inherit ANSI
	// from earlier in the string).
	return left + pad + overlay + "\x1b[0m" + right
}
