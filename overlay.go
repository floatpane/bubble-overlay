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
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
)

// Center places popup as a floating layer centered over base within a screen
// of screenW × screenH cells. The popup is a multi-line ANSI-styled string
// (e.g. the output of a lipgloss.Render call). Cells outside the popup region
// are taken from base; cells underneath the popup are replaced.
//
// This is the primary entry-point for modals, command palettes, and tooltips:
//
//	composited := overlay.Center(myView, myModal, termWidth, termHeight)
func Center(base, popup string, screenW, screenH int) string {
	lines := strings.Split(popup, "\n")

	popupW := 0
	for _, l := range lines {
		if w := ansi.StringWidth(l); w > popupW {
			popupW = w
		}
	}
	popupH := len(lines)

	col := (screenW - popupW) / 2
	if col < 0 {
		col = 0
	}
	row := (screenH - popupH) / 2
	if row < 0 {
		row = 0
	}

	return Block(base, lines, row, col)
}

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

	// Re-emit the SGR state that was active in the base at the point
	// just past the overlay, so styling continues correctly on the right
	// side. The \x1b[0m prevents overlay styles from bleeding into
	// surrounding cells; sgr restores the base's active styling.
	sgr := sgrStateAt(base, rightStart)

	return left + pad + overlay + "\x1b[0m" + sgr + right
}

// sgrStateAt returns the concatenation of all active SGR escape sequences
// at the given visible column position in s. This can be prepended to a
// substring of s to restore the styling that was active at that position.
func sgrStateAt(s string, col int) string {
	visible := 0
	var currentSGR strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' {
			j := i + 1
			if j < len(s) && s[j] == '[' {
				j++
				for j < len(s) {
					if s[j] >= 0x40 && s[j] <= 0x7e {
						j++
						break
					}
					j++
				}
			} else if j < len(s) {
				j++
			}
			seq := s[i:j]
			if len(seq) >= 3 && seq[0] == '\x1b' && seq[1] == '[' && seq[len(seq)-1] == 'm' {
				body := seq[2 : len(seq)-1]
				if body == "" || body == "0" {
					currentSGR.Reset()
				} else {
					currentSGR.WriteString(seq)
				}
			}
			i = j
			continue
		}
		if visible >= col {
			break
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		visible++
	}
	return currentSGR.String()
}

// BlockFloat paints block as a floating popup anchored near the cursor
// position (cursorRow, cursorCol) within a screen of screenW × screenH
// cells. The popup is placed below the cursor by default; if there is
// insufficient space below, it is placed above. The column is aligned to
// the cursor and shifted left if the popup would overflow the right edge.
//
// Unlike [Center], this is designed for autocomplete dropdowns, tooltips,
// and floating suggestion menus that appear at the cursor position.
func BlockFloat(base string, block []string, cursorRow, cursorCol, screenW, screenH int) string {
	if len(block) == 0 {
		return base
	}

	popupH := len(block)
	popupW := 0
	for _, l := range block {
		if w := ansi.StringWidth(l); w > popupW {
			popupW = w
		}
	}

	// Vertical: prefer below cursor, fall back to above
	row := cursorRow + 1
	if row+popupH > screenH {
		row = cursorRow - popupH
		if row < 0 {
			row = 0
		}
	}

	// Horizontal: align to cursor, shift left on overflow
	col := cursorCol
	if col+popupW > screenW {
		col = screenW - popupW
		if col < 0 {
			col = 0
		}
	}

	return Block(base, block, row, col)
}
