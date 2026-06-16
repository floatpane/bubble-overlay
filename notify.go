// Notification types for bubble-overlay: Error (red), Warning (yellow), Info (cyan).
//
// Each Notification is a self-contained Bubble Tea model. Embed one inside your
// parent model, forward messages to it in Update, and call PlaceOn in View.
//
// # Dismiss modes
//
// Three modes control how a notification closes:
//
//   - [DismissOnKey]      — user presses a configured key (default "q").
//   - [DismissAfterTimer] — closes automatically after [WithDuration]; a
//     progress bar is shown; no key dismissal.
//   - [DismissEither]     — closes after [WithDuration] OR when the key is
//     pressed; progress bar is shown.
//
// # Placement
//
// Position the notification anywhere by passing [WithPosition]:
//
//	n := overlay.NewError(
//	    overlay.WithTitle("Connection lost"),
//	    overlay.WithPosition(2, 4),
//	)
//
// In View(), call PlaceOn to composite the box onto your base view:
//
//	func (m Model) View() string {
//	    base := m.renderBase()
//	    return m.notification.PlaceOn(base)
//	}
//
// When Done returns true, PlaceOn is a no-op and View returns "".
package overlay

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Kind identifies the semantic type and default styling of a [Notification].
type Kind uint8

const (
	KindError   Kind = iota // red
	KindWarning             // yellow
	KindInfo                // cyan
)

// Default Unicode icon characters prepended to the notification header.
// Override per-notification with [WithIcon].
const (
	DefaultErrorIcon   = "✗"
	DefaultWarningIcon = "⚠"
	DefaultInfoIcon    = "ℹ"
)

// DismissMode controls how a [Notification] can be closed.
type DismissMode uint8

const (
	// DismissOnKey closes the notification when the user presses the key set
	// by [WithKey] (default "q"). No timer or progress bar is shown.
	DismissOnKey DismissMode = iota

	// DismissAfterTimer closes the notification automatically after the
	// duration set by [WithDuration] (default 5 s). No key dismissal is
	// available. A progress bar counts down the remaining time.
	DismissAfterTimer

	// DismissEither closes the notification after [WithDuration] OR when the
	// user presses the key set by [WithKey], whichever comes first. A
	// progress bar is shown.
	DismissEither
)

// notifyTickMsg is the internal timer tick sent every notifyTickInterval.
type notifyTickMsg time.Time

const notifyTickInterval = 100 * time.Millisecond

// Notification is a Bubble Tea model for a dismissible, styled notification
// box. Use [NewError], [NewWarning], or [NewInfo] to create one.
//
// Notification is a value type — all Bubble Tea methods return updated copies.
type Notification struct {
	kind        Kind
	title       string
	message     string
	icon        string
	dismissMode DismissMode
	key         string
	duration    time.Duration
	row         int
	col         int
	width       int
	elapsed     time.Duration
	done        bool
}

// Option configures a [Notification].
type Option func(*Notification)

// WithTitle sets the notification title shown next to the icon.
func WithTitle(title string) Option { return func(n *Notification) { n.title = title } }

// WithMessage sets the body text shown below the title.
func WithMessage(msg string) Option { return func(n *Notification) { n.message = msg } }

// WithKey sets the key string used to dismiss the notification (e.g. "q",
// "esc", "enter"). Defaults to "q".
func WithKey(key string) Option { return func(n *Notification) { n.key = key } }

// WithDuration sets the auto-close countdown for timer-based [DismissMode]
// values. Defaults to 5 s.
func WithDuration(d time.Duration) Option { return func(n *Notification) { n.duration = d } }

// WithDismissMode sets how the notification is closed. Defaults to
// [DismissOnKey].
func WithDismissMode(m DismissMode) Option { return func(n *Notification) { n.dismissMode = m } }

// WithPosition sets the (row, col) cell at which the notification is painted
// over the base view when [Notification.PlaceOn] is called. Both are
// 0-indexed. Defaults to (0, 0).
func WithPosition(row, col int) Option {
	return func(n *Notification) { n.row = row; n.col = col }
}

// WithIcon overrides the default Unicode icon character shown before the title.
func WithIcon(icon string) Option { return func(n *Notification) { n.icon = icon } }

// WithWidth sets the inner content width of the notification box in terminal
// cells. The rendered box is slightly wider due to border and padding.
// Defaults to 40.
func WithWidth(w int) Option { return func(n *Notification) { n.width = w } }

func newNotification(kind Kind, icon string, opts []Option) Notification {
	n := Notification{
		kind:        kind,
		icon:        icon,
		dismissMode: DismissOnKey,
		key:         "q",
		duration:    5 * time.Second,
		width:       40,
	}
	for _, o := range opts {
		o(&n)
	}
	return n
}

// NewError returns a red error notification. Default dismiss mode is
// [DismissOnKey] with key "q".
func NewError(opts ...Option) Notification {
	return newNotification(KindError, DefaultErrorIcon, opts)
}

// NewWarning returns a yellow warning notification. Default dismiss mode is
// [DismissOnKey] with key "q".
func NewWarning(opts ...Option) Notification {
	return newNotification(KindWarning, DefaultWarningIcon, opts)
}

// NewInfo returns a cyan informational notification. Default dismiss mode is
// [DismissOnKey] with key "q".
func NewInfo(opts ...Option) Notification {
	return newNotification(KindInfo, DefaultInfoIcon, opts)
}

// Init implements [tea.Model]. It starts the internal tick timer when the
// [DismissMode] involves a countdown.
func (n Notification) Init() tea.Cmd {
	if n.done {
		return nil
	}
	if n.dismissMode == DismissAfterTimer || n.dismissMode == DismissEither {
		return notifyTick()
	}
	return nil
}

// Update implements [tea.Model]. It handles key presses (for [DismissOnKey]
// and [DismissEither]) and timer ticks (for [DismissAfterTimer] and
// [DismissEither]).
func (n Notification) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if n.done {
		return n, nil
	}
	switch m := msg.(type) {
	case tea.KeyMsg:
		if (n.dismissMode == DismissOnKey || n.dismissMode == DismissEither) && m.String() == n.key {
			n.done = true
		}
	case notifyTickMsg:
		if n.dismissMode == DismissAfterTimer || n.dismissMode == DismissEither {
			n.elapsed += notifyTickInterval
			if n.elapsed >= n.duration {
				n.done = true
				return n, nil
			}
			return n, notifyTick()
		}
	}
	return n, nil
}

// Done reports whether the notification has been dismissed. When true,
// [Notification.View] returns "" and [Notification.PlaceOn] returns the base
// unchanged.
func (n Notification) Done() bool { return n.done }

// PlaceOn composites the notification box onto base at the configured
// (row, col) position using [Block]. When [Done] is true, base is returned
// unchanged.
func (n Notification) PlaceOn(base string) string {
	if n.done {
		return base
	}
	return Block(base, strings.Split(n.View(), "\n"), n.row, n.col)
}

// View implements [tea.Model]. Returns the fully rendered, bordered
// notification box. Returns "" when [Done] is true.
func (n Notification) View() string {
	if n.done {
		return ""
	}

	color := kindColor(n.kind)
	accent := lipgloss.NewStyle().Foreground(color).Bold(true)
	body := lipgloss.NewStyle().Foreground(color)

	header := n.icon + "  " + kindLabel(n.kind)
	if n.title != "" {
		header += ": " + n.title
	}

	var parts []string
	parts = append(parts, accent.Render(header))
	if n.message != "" {
		parts = append(parts, body.Render(n.message))
	}

	switch n.dismissMode {
	case DismissAfterTimer, DismissEither:
		barWidth := n.width / 2
		if barWidth < 5 {
			barWidth = 5
		}
		bar := notifyProgressBar(n.elapsed, n.duration, barWidth)
		remaining := (n.duration - n.elapsed).Round(100 * time.Millisecond)
		footer := fmt.Sprintf("%s  %s", bar, remaining)
		if n.dismissMode == DismissEither && n.key != "" {
			footer += fmt.Sprintf("  [%s] close", n.key)
		}
		parts = append(parts, body.Render(footer))
	case DismissOnKey:
		if n.key != "" {
			parts = append(parts, body.Faint(true).Render(fmt.Sprintf("[%s] close", n.key)))
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1).
		Width(n.width).
		Render(strings.Join(parts, "\n"))
}

func notifyTick() tea.Cmd {
	return tea.Tick(notifyTickInterval, func(t time.Time) tea.Msg {
		return notifyTickMsg(t)
	})
}

// notifyProgressBar renders a Unicode block progress bar of the given width
// showing elapsed/total. █ = elapsed, ░ = remaining.
func notifyProgressBar(elapsed, total time.Duration, width int) string {
	if width <= 0 {
		return ""
	}
	ratio := float64(elapsed) / float64(total)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(float64(width) * ratio)
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func kindColor(k Kind) lipgloss.Color {
	switch k {
	case KindError:
		return lipgloss.Color("9") // bright red
	case KindWarning:
		return lipgloss.Color("11") // bright yellow
	default:
		return lipgloss.Color("14") // bright cyan
	}
}

func kindLabel(k Kind) string {
	switch k {
	case KindError:
		return "Error"
	case KindWarning:
		return "Warning"
	default:
		return "Info"
	}
}
