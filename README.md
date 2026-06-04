<div align="center">

# bubble-overlay

**ANSI-aware overlay painter for Bubble Tea / lipgloss views.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/floatpane/bubble-overlay)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/floatpane/bubble-overlay.svg)](https://pkg.go.dev/github.com/floatpane/bubble-overlay)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/floatpane/bubble-overlay)](https://github.com/floatpane/bubble-overlay/releases)
[![CI](https://github.com/floatpane/bubble-overlay/actions/workflows/ci.yml/badge.svg)](https://github.com/floatpane/bubble-overlay/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

`bubble-overlay` paints rectangular blocks of styled text on top of an
existing ANSI-styled string at a given `(row, col)` cell position. It's the
missing primitive for modals, popups, tooltips, filepickers, and floating
panels in [Bubble Tea](https://github.com/charmbracelet/bubbletea) /
[lipgloss](https://github.com/charmbracelet/lipgloss) apps — render your
base view as a single string, render your popup as another, then composite.

## Features

- **SGR-safe.** Existing ANSI styles around the overlaid region are preserved; cells underneath are replaced. The overlay is terminated with `\x1b[0m` so its styles don't bleed into the row's tail.
- **Cell-accurate.** Uses `charmbracelet/x/ansi` for width — wide-character / emoji handling matches what the terminal actually renders, not byte counts.
- **Auto-grow.** Overlays that extend past the bottom of the base string append new lines instead of truncating.
- **Auto-pad.** Overlays past the right edge of a short row are padded with spaces, so a popup on row 3 column 40 still lands correctly when the base row is only 10 cells wide.
- **Tiny.** Two functions. No state. Drop-in.

## Install

```bash
go get github.com/floatpane/bubble-overlay
```

Requires Go 1.26+.

## Usage

```go
package main

import (
    "fmt"

    "github.com/charmbracelet/lipgloss"
    "github.com/floatpane/bubble-overlay"
)

func main() {
    base := lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")).
        Render("a quiet inbox view\nwith two lines\nand a third")

    popup := lipgloss.NewStyle().
        Background(lipgloss.Color("57")).
        Foreground(lipgloss.Color("231")).
        Padding(0, 1).
        Render("are you sure?\nyes / no")

    block := strings.Split(popup, "\n")
    fmt.Println(overlay.Block(base, block, 1, 4))
}
```

### Center a floating popup

For the common case of a centered modal — command palette, confirmation dialog,
tooltip — `Center` handles the positioning automatically:

```go
composited := overlay.Center(baseView, popupView, termWidth, termHeight)
```

It computes the centered `(row, col)` from the popup's visual size and your
screen dimensions, then calls `Block`. Rows/columns clamp to zero if the popup
is larger than the screen.

### Low-level API

```go
// Center places popup as a floating layer centered over base within a screen
// of screenW × screenH cells.
func Center(base, popup string, screenW, screenH int) string

// Paint a multi-line block on top of base at (row, col).
func Block(base string, block []string, row, col int) string

// Paint a single overlay line on top of base at col.
func Line(base, overlay string, col int) string
```

## When to use this

You have a Bubble Tea `View()` returning a styled multi-line string and
you want to render a modal/popup over it without:

- Re-rendering the base view with a "modal-shaped hole" cut out of it.
- Walking the ANSI sequences yourself.
- Truncating styles that span across the modal region.

`bubble-overlay` does the composite for you. The base view stays a single
string; the modal stays a single string; you call `Block` and emit the
result.

## Documentation

Full API reference: [pkg.go.dev/github.com/floatpane/bubble-overlay](https://pkg.go.dev/github.com/floatpane/bubble-overlay)

## Contributing

PRs welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

Report vulnerabilities privately via [SECURITY.md](SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
