package overlay

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

// --- Constructor defaults ---

func TestNewErrorDefaults(t *testing.T) {
	n := NewError()
	if n.kind != KindError {
		t.Fatal("kind should be KindError")
	}
	if n.icon != DefaultErrorIcon {
		t.Fatalf("icon: got %q, want %q", n.icon, DefaultErrorIcon)
	}
	if n.dismissMode != DismissOnKey {
		t.Fatal("default dismissMode should be DismissOnKey")
	}
	if n.key != "q" {
		t.Fatalf("default key: got %q, want %q", n.key, "q")
	}
	if n.duration != 5*time.Second {
		t.Fatalf("default duration: got %v, want 5s", n.duration)
	}
	if n.width != 40 {
		t.Fatalf("default width: got %d, want 40", n.width)
	}
	if n.Done() {
		t.Fatal("new notification should not be done")
	}
}

func TestNewWarningDefaults(t *testing.T) {
	n := NewWarning()
	if n.kind != KindWarning {
		t.Fatal("kind should be KindWarning")
	}
	if n.icon != DefaultWarningIcon {
		t.Fatalf("icon: got %q, want %q", n.icon, DefaultWarningIcon)
	}
}

func TestNewInfoDefaults(t *testing.T) {
	n := NewInfo()
	if n.kind != KindInfo {
		t.Fatal("kind should be KindInfo")
	}
	if n.icon != DefaultInfoIcon {
		t.Fatalf("icon: got %q, want %q", n.icon, DefaultInfoIcon)
	}
}

// --- Options ---

func TestOptionsApplied(t *testing.T) {
	n := NewError(
		WithTitle("oops"),
		WithMessage("something failed"),
		WithKey("esc"),
		WithDuration(3*time.Second),
		WithPosition(5, 10),
		WithWidth(50),
		WithIcon("!"),
		WithDismissMode(DismissEither),
	)
	if n.title != "oops" {
		t.Fatalf("title: got %q", n.title)
	}
	if n.message != "something failed" {
		t.Fatalf("message: got %q", n.message)
	}
	if n.key != "esc" {
		t.Fatalf("key: got %q", n.key)
	}
	if n.duration != 3*time.Second {
		t.Fatalf("duration: got %v", n.duration)
	}
	if n.row != 5 || n.col != 10 {
		t.Fatalf("position: got (%d,%d)", n.row, n.col)
	}
	if n.width != 50 {
		t.Fatalf("width: got %d", n.width)
	}
	if n.icon != "!" {
		t.Fatalf("icon: got %q", n.icon)
	}
	if n.dismissMode != DismissEither {
		t.Fatal("dismissMode should be DismissEither")
	}
}

// --- DismissOnKey ---

func TestDismissOnKeyClosesOnCorrectKey(t *testing.T) {
	n := NewError(WithKey("q"))
	m, _ := n.Update(tea.KeyPressMsg{Text: "q"})
	if !m.(Notification).Done() {
		t.Fatal("should be done after correct key press")
	}
}

func TestDismissOnKeyIgnoresWrongKey(t *testing.T) {
	n := NewError(WithKey("q"))
	m, _ := n.Update(tea.KeyPressMsg{Text: "x"})
	if m.(Notification).Done() {
		t.Fatal("should not be done on wrong key")
	}
}

func TestDismissOnKeyIgnoresTick(t *testing.T) {
	n := NewError(WithKey("q"))
	m, _ := n.Update(notifyTickMsg(time.Time{}))
	if m.(Notification).Done() {
		t.Fatal("key-only mode should not close on tick")
	}
}

// --- DismissAfterTimer ---

func TestDismissAfterTimerIgnoresKey(t *testing.T) {
	n := NewError(WithDismissMode(DismissAfterTimer), WithKey("q"))
	m, _ := n.Update(tea.KeyPressMsg{Text: "q"})
	if m.(Notification).Done() {
		t.Fatal("timer-only mode must not close on key press")
	}
}

func TestDismissAfterTimerAdvancesElapsed(t *testing.T) {
	n := NewError(WithDismissMode(DismissAfterTimer), WithDuration(500*time.Millisecond))
	m, cmd := n.Update(notifyTickMsg(time.Time{}))
	n2 := m.(Notification)
	if n2.elapsed != notifyTickInterval {
		t.Fatalf("elapsed: got %v, want %v", n2.elapsed, notifyTickInterval)
	}
	if cmd == nil {
		t.Fatal("should return another tick cmd while not done")
	}
}

func TestDismissAfterTimerExpiresWhenElapsedReachesDuration(t *testing.T) {
	n := NewError(WithDismissMode(DismissAfterTimer), WithDuration(200*time.Millisecond))
	n.elapsed = 100 * time.Millisecond // one tick away from expiry
	m, _ := n.Update(notifyTickMsg(time.Time{}))
	if !m.(Notification).Done() {
		t.Fatal("should be done once elapsed >= duration")
	}
}

// --- DismissEither ---

func TestDismissEitherClosesOnKey(t *testing.T) {
	n := NewError(WithDismissMode(DismissEither), WithKey("esc"))
	m, _ := n.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if !m.(Notification).Done() {
		t.Fatal("DismissEither should close on key press")
	}
}

func TestDismissEitherClosesOnTimer(t *testing.T) {
	n := NewError(WithDismissMode(DismissEither), WithDuration(100*time.Millisecond))
	m, _ := n.Update(notifyTickMsg(time.Time{}))
	if !m.(Notification).Done() {
		t.Fatal("DismissEither should close when timer expires")
	}
}

// --- Init ---

func TestInitReturnsTickForTimerModes(t *testing.T) {
	for _, mode := range []DismissMode{DismissAfterTimer, DismissEither} {
		n := NewError(WithDismissMode(mode))
		if n.Init() == nil {
			t.Fatalf("Init should return tick cmd for mode %v", mode)
		}
	}
}

func TestInitReturnsNilForKeyOnly(t *testing.T) {
	n := NewError(WithDismissMode(DismissOnKey))
	if n.Init() != nil {
		t.Fatal("Init should return nil for DismissOnKey")
	}
}

func TestInitReturnsNilWhenDone(t *testing.T) {
	n := NewError(WithDismissMode(DismissAfterTimer))
	n.done = true
	if n.Init() != nil {
		t.Fatal("Init should return nil when already done")
	}
}

// --- View ---

func TestViewContainsTitleAndMessage(t *testing.T) {
	n := NewError(WithTitle("Disk full"), WithMessage("Clean up space"))
	v := n.Render()
	if !strings.Contains(v, "Disk full") {
		t.Fatalf("view missing title: %q", v)
	}
	if !strings.Contains(v, "Clean up space") {
		t.Fatalf("view missing message: %q", v)
	}
}

func TestViewContainsIcon(t *testing.T) {
	n := NewError()
	if !strings.Contains(n.Render(), DefaultErrorIcon) {
		t.Fatal("view missing default error icon")
	}
	n2 := NewWarning()
	if !strings.Contains(n2.Render(), DefaultWarningIcon) {
		t.Fatal("view missing default warning icon")
	}
	n3 := NewInfo()
	if !strings.Contains(n3.Render(), DefaultInfoIcon) {
		t.Fatal("view missing default info icon")
	}
}

func TestViewContainsKeyHintForDismissOnKey(t *testing.T) {
	n := NewError(WithKey("esc"))
	if !strings.Contains(n.Render(), "[esc]") {
		t.Fatal("key-only view should show key hint")
	}
}

func TestViewContainsProgressBarForTimerModes(t *testing.T) {
	for _, mode := range []DismissMode{DismissAfterTimer, DismissEither} {
		n := NewError(WithDismissMode(mode))
		v := n.Render()
		if !strings.Contains(v, "░") && !strings.Contains(v, "█") {
			t.Fatalf("timer mode %v should show progress bar in view", mode)
		}
	}
}

func TestViewEmptyWhenDone(t *testing.T) {
	n := NewError()
	m, _ := n.Update(tea.KeyPressMsg{Text: "q"})
	if v := m.(Notification).Render(); v != "" {
		t.Fatalf("View should be empty when done, got %q", v)
	}
}

func TestViewCustomIcon(t *testing.T) {
	n := NewInfo(WithIcon("→"))
	if !strings.Contains(n.Render(), "→") {
		t.Fatal("view should use custom icon")
	}
	if strings.Contains(n.Render(), DefaultInfoIcon) {
		t.Fatal("view should not contain default icon when overridden")
	}
}

// --- PlaceOn ---

func TestPlaceOnReturnBaseWhenDone(t *testing.T) {
	n := NewError(WithPosition(0, 0))
	m, _ := n.Update(tea.KeyPressMsg{Text: "q"})
	base := "hello world"
	if got := m.(Notification).PlaceOn(base); got != base {
		t.Fatalf("PlaceOn should return base unchanged when done, got %q", got)
	}
}

func TestPlaceOnCompositsWhenNotDone(t *testing.T) {
	base := strings.Repeat(strings.Repeat(".", 60)+"\n", 10)
	n := NewError(WithTitle("Oops"), WithPosition(2, 5))
	result := n.PlaceOn(base)
	if result == base {
		t.Fatal("PlaceOn should modify base when not done")
	}
	if !strings.Contains(result, "Oops") {
		t.Fatal("PlaceOn result should contain notification content")
	}
}

// --- Progress bar ---

func TestProgressBarEmpty(t *testing.T) {
	got := notifyProgressBar(0, 10*time.Second, 10)
	if got != "░░░░░░░░░░" {
		t.Fatalf("empty bar: got %q", got)
	}
}

func TestProgressBarFull(t *testing.T) {
	got := notifyProgressBar(10*time.Second, 10*time.Second, 10)
	if got != "██████████" {
		t.Fatalf("full bar: got %q", got)
	}
}

func TestProgressBarHalf(t *testing.T) {
	got := notifyProgressBar(5*time.Second, 10*time.Second, 10)
	if got != "█████░░░░░" {
		t.Fatalf("half bar: got %q", got)
	}
}

func TestProgressBarZeroWidth(t *testing.T) {
	got := notifyProgressBar(5*time.Second, 10*time.Second, 0)
	if got != "" {
		t.Fatalf("zero-width bar should be empty, got %q", got)
	}
}

// --- No mutation of original ---

func TestUpdateDoesNotMutateOriginal(t *testing.T) {
	n := NewError(WithDismissMode(DismissAfterTimer), WithDuration(time.Second))
	n.Update(notifyTickMsg(time.Time{})) //nolint
	if n.elapsed != 0 {
		t.Fatal("Update must not mutate the original value")
	}
}
